package server

import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/internal/mcp/proxy"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
	serverInfo, ok := m.cache.LoadMcpServer(cache.NewMcpValue)
	if !ok {
		m.log.Errorf("LoadMcpServer error: %v", "加载内存serverInfo缓存信息失败")
		return
	}
	if serverInfo == nil || serverInfo.Tools == nil || len(serverInfo.Tools) == 0 {
		m.log.Errorf("LoadMcpServer error,serverInfo:%+v", serverInfo)
		return
	}
	for _, tool := range serverInfo.Tools {
		var toolSchema protocol.InputSchema
		err := json.Unmarshal([]byte(tool.ToolSchema), &toolSchema)
		if err != nil {
			m.log.Errorf("Failed to unmarshal tool schema: %v", err)
			continue
		}
		name := tool.Name
		if tool.IsRepeat == _const.StatusDisplay {
			name = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
		}
		toolInfo := &protocol.Tool{
			Name:           name,
			Description:    tool.Description,
			InputSchema:    toolSchema,
			OutputSchema:   protocol.OutputSchema{},
			Annotations:    nil,
			RawInputSchema: nil,
		}
		m.Server.RegisterTool(toolInfo, m.httpProxy.HandleHttpProxy)
	}
	m.cache.ClearCache(cache.OldMcpValue)

	return
}

func (m *McpServerManager) UnRegisterToolFromCache() {
	serverInfo, ok := m.cache.LoadMcpServer(cache.OldMcpValue)
	if !ok {
		m.log.Errorf("UnRegisterToolFromCache error: %v", "加载内存Old cache serverInfo缓存信息失败")
		return
	}
	if serverInfo == nil || serverInfo.Tools == nil || len(serverInfo.Tools) == 0 {
		return
	}
	for _, tool := range serverInfo.Tools {
		name := tool.Name
		if tool.IsRepeat == _const.StatusDisplay {
			name = tool.Name + "_" + strconv.Itoa(int(tool.McpServerId)) + tool.SerialNumber
		}
		m.Server.UnregisterTool(name)
	}
	return
}
