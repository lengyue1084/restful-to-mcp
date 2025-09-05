package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHttpProxy,
)

type HttpProxy struct {
	log        *logger.Logger
	cache      *cache.MemoryCache
	httpClient *http.Client
}

func NewHttpProxy(log *logger.Logger, cache *cache.MemoryCache) *HttpProxy {
	return &HttpProxy{
		log:   log,
		cache: cache,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ToolMetadata 工具元数据
type ToolMetadata struct {
	Method      string            `json:"method"`
	Endpoint    string            `json:"endpoint"`
	Urls        []string          `json:"urls"`
	Headers     map[string]string `json:"headers"`
	ContentType string            `json:"contentType"`
}

// RequestParams 请求参数结构
type RequestParams struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    interface{}
}

func (h *HttpProxy) HandleHttpProxy(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 从上下文获取服务令牌
	if serviceToken, ok := ctx.Value(_const.ServiceToken).(string); ok {
		h.log.Infof("Service Token: %s", serviceToken)
	}

	// 查找工具信息
	toolInfo, err := h.findToolInfo(req.Name)
	if err != nil {
		return nil, fmt.Errorf("查找工具信息失败: %w", err)
	}

	// 获取工具参数配置
	argsSlice, err := h.getToolArgs(toolInfo)
	if err != nil {
		return nil, fmt.Errorf("获取工具参数配置失败: %w", err)
	}

	// 获取工具元数据
	toolMetadata, err := h.getToolMetadata(toolInfo)
	if err != nil {
		return nil, fmt.Errorf("获取工具元数据失败: %w", err)
	}

	// 构建HTTP请求参数
	requestParams, err := h.buildRequestParams(req.Arguments, argsSlice, toolMetadata)
	if err != nil {
		return nil, fmt.Errorf("构建请求参数失败: %w", err)
	}

	// 发送HTTP请求
	response, err := h.sendHttpRequest(ctx, requestParams)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}

	return &protocol.CallToolResult{
		IsError: false,
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(response),
			},
		},
	}, nil
}

// findToolInfo 查找工具信息
func (h *HttpProxy) findToolInfo(toolName string) (*model.McpTools, error) {
	mcpServerList, ok := h.cache.LoadMcpServer(cache.NewMcpValue)
	if !ok {
		return nil, fmt.Errorf("加载内存serverInfo缓存信息失败")
	}
	for _, mcpServer := range mcpServerList {
		if len(mcpServer.Tools) > 0 {
			for _, tool := range mcpServer.Tools {
				// 处理重复工具名称
				actualToolName := tool.Name
				if tool.IsRepeat == _const.CommonStatusYes && tool.SerialNumber != "" {
					actualToolName = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
				}

				if toolName == actualToolName {
					if tool.McpServerType != _const.McpServerTypeOpenapi {
						return nil, fmt.Errorf("该工具只支持%s类型", _const.McpServerTypeOpenapi)
					}
					var urls []string
					err := json.Unmarshal([]byte(mcpServer.Urls), &urls)
					if err != nil {
						return nil, err
					}
					var tmpEndpoint string
					//处理多个url的问题
					selectedURLSlice := h.selectValidURL(urls)
					for _, urlVal := range selectedURLSlice {
						if strings.Contains(tool.Endpoint, "{{.Config.url}}") {
							tmpEndpoint += strings.ReplaceAll(tool.Endpoint, "{{.Config.url}}", urlVal) + "|"
						}
					}

					tmpEndpoint = strings.TrimSuffix(tmpEndpoint, "|")
					tool.Endpoint = tmpEndpoint
					return tool, nil
				}

			}
		}
	}

	return nil, fmt.Errorf("未找到工具: %s", toolName)
}

// selectValidURL 选择第一个有效的HTTP/HTTPS URL
func (h *HttpProxy) selectValidURL(urls []string) (urlsSlice []string) {
	for _, urlVal := range urls {
		// 检查是否为有效的HTTP/HTTPS URL
		if strings.HasPrefix(urlVal, "http://") || strings.HasPrefix(urlVal, "https://") {
			urlsSlice = append(urlsSlice, urlVal)
		}
	}
	return urlsSlice
}

// getToolArgs 获取工具参数配置
func (h *HttpProxy) getToolArgs(toolInfo *model.McpTools) ([]*config.ArgConfig, error) {
	var argsSlice []*config.ArgConfig
	err := json.Unmarshal([]byte(toolInfo.Args), &argsSlice)
	if err != nil {
		h.log.Errorf("解析工具参数失败: %v", err)
		return nil, fmt.Errorf("解析工具参数失败: %w", err)
	}
	return argsSlice, nil
}

// buildToolMetadataFromFields 从字段构建工具元数据（备用方案）
func (h *HttpProxy) getToolMetadata(toolInfo *model.McpTools) (*ToolMetadata, error) {
	if toolInfo.Endpoint == "" {
		h.log.Errorf("工具URL为空")
		return &ToolMetadata{}, fmt.Errorf("工具URL为空")
	}
	var urls = strings.Split(toolInfo.Endpoint, "|")
	if len(urls) == 0 {
		h.log.Errorf("工具URL为空")
		return &ToolMetadata{}, fmt.Errorf("工具URL为空")
	}
	metadata := &ToolMetadata{
		Method:      toolInfo.Method,
		Endpoint:    urls[0],
		Headers:     make(map[string]string),
		Urls:        urls,
		ContentType: "application/json", //目前默认为JSON
	}

	// 解析请求头
	if toolInfo.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(toolInfo.Headers), &headers); err != nil {
			h.log.Errorf("解析请求头失败: %v", err)
			return &ToolMetadata{}, fmt.Errorf("解析请求头失败: %w", err)
		}
		metadata.Headers = headers
	}

	// 解析URLs
	// 这里假设URL存储在某个地方，需要根据实际情况调整
	return metadata, nil
}

// buildRequestParams 构建请求参数
func (h *HttpProxy) buildRequestParams(inputArgs map[string]interface{}, argsSlice []*config.ArgConfig, toolMetadata *ToolMetadata) (*RequestParams, error) {
	params := &RequestParams{
		URL:     toolMetadata.Endpoint,
		Method:  toolMetadata.Method,
		Headers: make(map[string]string),
	}

	// 复制工具配置的请求头
	for key, value := range toolMetadata.Headers {
		params.Headers[key] = value
	}

	// 分类参数
	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headerParams := make(map[string]string)
	bodyParams := make(map[string]interface{})

	// 根据参数位置分类
	for _, argConfig := range argsSlice {
		if value, exists := inputArgs[argConfig.Name]; exists {
			switch argConfig.Position {
			case "path":
				pathParams[argConfig.Name] = fmt.Sprintf("%v", value)
			case "query":
				h.addQueryParam(queryParams, argConfig, value)
			case "header":
				headerParams[argConfig.Name] = fmt.Sprintf("%v", value)
			case "body":
				bodyParams[argConfig.Name] = value
			}
		} else if argConfig.Required {
			return nil, fmt.Errorf("缺少必需参数: %s", argConfig.Name)
		}
	}

	// 构建URL
	baseURL := toolMetadata.Endpoint
	if baseURL == "" { //获取的默认为第一个
		return nil, fmt.Errorf("工具URL为空")
	}

	// 替换路径参数

	for name, value := range pathParams {
		placeholder := "{{.Args." + name + "}}"
		toolMetadata.Endpoint = strings.ReplaceAll(toolMetadata.Endpoint, placeholder, url.PathEscape(value))
		for key, urlVal := range toolMetadata.Urls {
			toolMetadata.Urls[key] = strings.ReplaceAll(urlVal, "{{.Args."+name+"}}", url.PathEscape(value))
		}
	}
	endpoint := toolMetadata.Endpoint

	// 构建完整URL
	fullURL := endpoint
	// 添加查询参数
	if len(queryParams) > 0 {
		queryStr := queryParams.Encode()
		if strings.Contains(fullURL, "?") {
			fullURL += "&" + queryStr
		} else {
			fullURL += "?" + queryStr
		}
	}
	params.URL = fullURL

	// 设置请求头参数
	for key, value := range headerParams {
		params.Headers[key] = value
	}

	// 设置请求体
	if len(bodyParams) > 0 {
		params.Body = bodyParams
		// 设置默认Content-Type
		if _, exists := params.Headers["Content-Type"]; !exists {
			if toolMetadata.ContentType != "" {
				params.Headers["Content-Type"] = toolMetadata.ContentType
			} else {
				params.Headers["Content-Type"] = "application/json"
			}
		}
	}

	return params, nil
}

// addQueryParam 添加查询参数
func (h *HttpProxy) addQueryParam(queryParams url.Values, argConfig *config.ArgConfig, value interface{}) {
	switch argConfig.Type {
	case "array":
		// 处理数组类型的查询参数
		if arr, ok := value.([]interface{}); ok {
			for _, item := range arr {
				queryParams.Add(argConfig.Name, fmt.Sprintf("%v", item))
			}
		} else {
			queryParams.Add(argConfig.Name, fmt.Sprintf("%v", value))
		}
	default:
		queryParams.Add(argConfig.Name, fmt.Sprintf("%v", value))
	}
}

// sendHttpRequest 发送HTTP请求
func (h *HttpProxy) sendHttpRequest(ctx context.Context, params *RequestParams) ([]byte, error) {
	// 构建请求体
	var body io.Reader
	if params.Body != nil {
		bodyBytes, err := json.Marshal(params.Body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, params.Method, params.URL, body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range params.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发送请求
	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
