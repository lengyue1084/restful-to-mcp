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
	"flow-bridge-mcp/pkg/tool"
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
	serviceToken, ok := ctx.Value(_const.ServiceToken).(string)
	if ok {
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
		return nil, fmt.Errorf("获取工具参数配置失败: %+v", err)
	}

	// 获取工具元数据
	toolMetadata, err := h.getToolMetadata(toolInfo)
	if err != nil {
		return nil, fmt.Errorf("获取工具元数据失败: %+v", err)
	}
	var security *config.Security
	if toolInfo.IsAuth == _const.IsAuthYes && toolInfo.Security != "" {
		err := json.Unmarshal([]byte(toolInfo.Security), &security)
		if err != nil {
			return nil, fmt.Errorf("解析工具安全配置失败: %+v", err)
		}
	}

	// 构建HTTP请求参数
	requestParams, err := h.buildRequestParams(req.Arguments, argsSlice, toolMetadata, security, serviceToken)
	if err != nil {
		return nil, fmt.Errorf("构建请求参数失败: %+v", err)
	}

	var (
		//重试次数
		retryCount = 0
		//最大重试次数
		maxRetries = 3
		response   []byte
	)
	for retryCount < maxRetries {
		response, err = h.sendHttpRequest(ctx, requestParams)
		if err == nil {
			break
		}
		retryCount++
		h.log.Warnf("请求%s失败，重试次数: %d, 错误: %v", requestParams.URL, retryCount, err)

		// 如果达到最大重试次数，尝试备用URL
		if retryCount >= maxRetries && len(toolMetadata.Urls) > 1 {
			h.log.Infof("主URL重试失败，尝试备用URL")
			// 保存原始URL用于恢复
			originalURL := requestParams.URL
			// 使用备用URL
			requestParams.URL = toolMetadata.Urls[1]
			response, err = h.sendHttpRequest(ctx, requestParams)
			if err != nil {
				// 备用URL也失败，恢复原始URL并返回错误
				requestParams.URL = originalURL
				h.log.Errorf("所有URL:%s尝试失败，主URL错误: %v, 备用URL错误: %v", requestParams.URL, err, err)
				return nil, fmt.Errorf("所有URL尝试失败，主URL错误: %w, 备用URL错误: %v", err, err)
			}
			// 备用URL成功
			break
		}
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
		h.log.Errorf("LoadMcpServer error: %v", "加载内存serverInfo缓存信息失败")
		return nil, fmt.Errorf("加载内存serverInfo缓存信息失败")
	}
	var toolInfo model.McpTools
	for _, mcpServer := range mcpServerList {
		if len(mcpServer.Tools) > 0 {
			for _, tool := range mcpServer.Tools {
				toolInfo = *tool
				// 处理重复工具名称
				actualToolName := tool.Name
				if tool.IsRepeat == _const.CommonStatusYes && tool.SerialNumber != "" {
					actualToolName = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
				}

				if toolName == actualToolName {
					if tool.McpServerType != _const.McpServerTypeOpenapi {
						h.log.Errorf("该工具只支持%s类型", _const.McpServerTypeOpenapi)
						return nil, fmt.Errorf("该工具只支持%s类型", _const.McpServerTypeOpenapi)
					}
					var urls []string
					err := json.Unmarshal([]byte(mcpServer.Urls), &urls)
					if err != nil {
						h.log.Errorf("mcpServerInfo.Urls json转换错误，err:%+v", err)
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
					toolInfo.Endpoint = tmpEndpoint
					return &toolInfo, nil
				}

			}
		}
	}
	h.log.Errorf("未找到工具: %s", toolName)
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
func (h *HttpProxy) buildRequestParams(inputArgs map[string]interface{}, argsSlice []*config.ArgConfig, toolMetadata *ToolMetadata, security *config.Security, serviceToken string) (*RequestParams, error) {
	params := &RequestParams{
		URL:     toolMetadata.Endpoint,
		Method:  toolMetadata.Method,
		Headers: make(map[string]string),
	}

	// 复制工具配置的请求头
	for key, value := range toolMetadata.Headers {
		params.Headers[key] = value
	}

	// 处理鉴权参数
	pathParams := make(map[string]string)
	queryParams := make(url.Values)
	headerParams := make(map[string]string)
	bodyParams := make(map[string]interface{})
	if security != nil {
		switch security.Mode {
		case config.AuthModeApiKey:
			switch security.In {
			case config.AuthPositionHeader:
				params.Headers[security.Name] = serviceToken
			case config.AuthPositionQuery:
				queryParams.Add(security.Name, serviceToken)
			default:
				return nil, fmt.Errorf("不支持的API密钥位置: %s", security.In)
			}
		case config.AuthModeHttp:
			params.Headers["Authorization"] = fmt.Sprintf("%s %s", tool.Capitalize(security.Scheme), serviceToken)
		default:
			return nil, fmt.Errorf("不支持的鉴权模式: %s", security.Mode)
		}
	}

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
	if resp.StatusCode != 200 {
		h.log.Errorf("HTTP请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
		return respBody, fmt.Errorf("HTTP请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
