package server

import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/proxy"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"log"
	"strconv"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewMcpServerManager,
	NewMcpTransport,
)

type McpTransportType string

const (
	TypeSSE        McpTransportType = "sse"
	TypeStreamable McpTransportType = "streamable-http"
)

// McpServerManager 暂时只支持streamable
type McpServerManager struct {
	Server                  *server.Server
	log                     *logger.Logger
	streamableHttpTransprot *StreamableHttpTransprot
	cache                   *cache.MemoryCache
	httpProxy               *proxy.HttpProxy
}

func NewMcpServerManager(streamableHttpTransprot *StreamableHttpTransprot, httpProxy *proxy.HttpProxy, log *logger.Logger, cache *cache.MemoryCache) (mcpServerManager *McpServerManager, err error) {
	// 创建 streamable server
	streamableServer, err := server.NewServer(streamableHttpTransprot.StreamableTransport)
	if err != nil {
		return nil, err
	}
	return &McpServerManager{
		Server:                  streamableServer,
		log:                     log,
		streamableHttpTransprot: streamableHttpTransprot,
		cache:                   cache,
		httpProxy:               httpProxy,
	}, err
}
func (m *McpServerManager) Run(ctx context.Context) error {
	// 启动 MCP 服务器
	//m.Server.Use(m.PanicRecoveryMiddleware(), m.ServerToolListByServerIdMiddleware()) //动态加载工具，全局的中间件会失效
	m.Server.SetToolFilter(m.FilterToolsByServer)
	serverErrChan := make(chan error, 1)
	go func() {
		m.log.WithContext(ctx).Info("Starting MCP server")
		serverErrChan <- m.Server.Run()
	}()
	select {
	case err := <-serverErrChan:
		m.log.WithContext(ctx).Error("MCP server error: %v", err)
		return err
	case <-ctx.Done():
		m.log.WithContext(ctx).Info("Shutting down MCP server")
		m.Server.Shutdown(context.Background())
		return ctx.Err()
	}

}
func (m *McpServerManager) FilterToolsByServer(ctx context.Context, tools []*protocol.Tool) []*protocol.Tool {

	filterTools := make([]*protocol.Tool, 0)

	if value, ok := ctx.Value(_const.ServerToken).(string); ok {
		if mcpServerList, ok := m.cache.LoadMcpServer(cache.NewMcpValue); ok {
			var mcpServerTools []*model.McpTools
			for _, mcpServer := range mcpServerList {
				if mcpServer.UUID == value {
					mcpServerTools = mcpServer.Tools
					break
				}
			}
			if len(mcpServerTools) > 0 {
				for _, tool := range mcpServerTools {
					if tool.IsShow == _const.StatusDisplay {
						toolName := tool.Name
						if tool.IsRepeat == _const.CommonStatusYes {
							toolName = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
						}
						for _, toolIem := range tools {
							if toolIem.Name == toolName {
								filterTools = append(filterTools, toolIem)
							}
						}
					}
				}
			}
		}
	}

	return filterTools
}

func (m *McpServerManager) ServerToolListByServerIdMiddleware() server.ToolMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {

			m.log.WithContext(ctx).Info("HandleHttpProxy: %v", req)

			return next(ctx, req)
		}
	}
}

// PanicRecoveryMiddleware returns a panic recovery middleware
func (m *McpServerManager) PanicRecoveryMiddleware() server.ToolMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Middleware] Recovered from panic in tool %s: %v", req.Name, r)
				}
			}()

			return next(ctx, req)
		}
	}
}

// HandleConnection 专门处理HTTP连接
func (m *McpServerManager) HandleConnection(c *gin.Context) {
	// 处理 MCP 连接
	// 在任何处理之前设置所有必需的响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache, no-transform")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用代理缓冲
	c.Header("Access-Control-Allow-Origin", "*")

	// 立即刷新响应头
	c.Writer.WriteHeaderNow()
	m.streamableHttpTransprot.StreamableHandler.HandleMCP().ServeHTTP(c.Writer, c.Request)
}

func (m *McpServerManager) RegisterToolFromCache() {
	m.UnRegisterToolFromCache()
	serverInfoList, ok := m.cache.LoadMcpServer(cache.NewMcpValue)
	if !ok {
		m.log.Errorf("LoadMcpServer error: %v", "加载内存serverInfo缓存信息失败")
		return
	}
	if len(serverInfoList) == 0 {
		m.log.Errorf("LoadMcpServer error: %v", "没有获取到serverInfoList信息")
		return
	}
	for _, serverInfo := range serverInfoList {
		if serverInfo == nil || serverInfo.Tools == nil || len(serverInfo.Tools) == 0 {
			m.log.Errorf("LoadMcpServer error,serverInfo:%+v", serverInfo)
			return
		}
		for _, tool := range serverInfo.Tools {
			if tool.IsShow != _const.StatusDisplay {
				continue
			}
			var toolSchema protocol.InputSchema
			err := json.Unmarshal([]byte(tool.ToolSchema), &toolSchema)
			if err != nil {
				m.log.Errorf("Failed to unmarshal tool schema: %v", err)
				continue
			}
			name := tool.Name
			if tool.IsRepeat == _const.CommonStatusYes {
				name = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
			}
			//m.Server.Use(m.PanicRecoveryMiddleware(), m.ServerToolListByServerIdMiddleware()) //动态加载工具，全局的中间件会失效
			toolInfo := &protocol.Tool{
				Name:           name,
				Description:    tool.Description,
				InputSchema:    toolSchema,
				OutputSchema:   protocol.OutputSchema{},
				Annotations:    nil,
				RawInputSchema: nil,
			}
			authentication := m.authenticationMiddleware()
			// 创建带上下文信息的处理函数
			handler := m.createContextAwareHandler()
			m.Server.RegisterTool(toolInfo, handler, authentication)
		}
	}

	m.cache.ClearCache(cache.OldMcpValue)
	return
}

func (m *McpServerManager) authenticationMiddleware() server.ToolMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
			platformToken, ok := ctx.Value(_const.PlatformToken).(string)
			if !ok {
				return nil, fmt.Errorf("无效的%s", _const.PlatformToken)
			}
			serviceToken, ok := ctx.Value(_const.ServiceToken).(string)
			if !ok {
				return nil, fmt.Errorf("无效的%s", _const.ServiceToken)
			}
			mcpServerList, ok := m.cache.LoadMcpServer(cache.NewMcpValue)
			if !ok {
				return nil, fmt.Errorf("LoadMcpServer error: %v", "加载内存serverInfo缓存信息失败")
			}
			for _, serverInfo := range mcpServerList {
				if len(serverInfo.Tools) > 0 {
					for _, tool := range serverInfo.Tools {
						toolName := tool.Name
						if tool.IsRepeat == _const.CommonStatusYes && tool.SerialNumber != "" {
							toolName = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
						}
						if req.Name == toolName {
							if tool.IsPlatformAuth == _const.IsAuthYes && serverInfo.PlatformToken != platformToken {
								m.log.WithContext(ctx).Errorf("授权平台token%s:%s无效,方法:%s", _const.PlatformToken, platformToken, req.Name)
								return nil, fmt.Errorf("授权平台token%s无效", _const.PlatformToken)
							}
							if tool.IsAuth == _const.IsAuthYes && serverInfo.ServiceToken != serviceToken {
								m.log.WithContext(ctx).Errorf("授权服务token%s:%s无效,方法:%s", _const.ServiceToken, serviceToken, req.Name)
								return nil, fmt.Errorf("授权服务token%s无效", _const.ServiceToken)
							}
							if tool.IsShow == _const.StatusHidden {
								m.log.WithContext(ctx).Errorf("隐藏方法%s", req.Name)
								return nil, fmt.Errorf("该方法%s不可用,请核对后再试", req.Name)
							}
							break
						}
					}
				}
			}
			return next(ctx, req)
		}
	}
}

// 创建带上下文信息的处理函数
func (m *McpServerManager) createContextAwareHandler() func(context.Context, *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	return func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		return m.httpProxy.HandleHttpProxy(ctx, req)
	}
}

func (m *McpServerManager) UnRegisterToolFromCache() {
	serverInfoList, ok := m.cache.LoadMcpServer(cache.OldMcpValue)
	if !ok {
		m.log.Errorf("UnRegisterToolFromCache error: %v", "加载内存Old cache serverInfoList缓存信息失败")
		return
	}
	if len(serverInfoList) == 0 {
		m.log.Errorf("UnRegisterToolFromCache error: %v", "没有获取到serverInfoList信息")
		return
	}
	for _, serverInfo := range serverInfoList {
		if serverInfo == nil || serverInfo.Tools == nil || len(serverInfo.Tools) == 0 {
			return
		}
		for _, tool := range serverInfo.Tools {
			name := tool.Name
			if tool.IsRepeat == _const.CommonStatusYes {
				name = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
			}
			m.Server.UnregisterTool(name)
		}
	}
	return
}
