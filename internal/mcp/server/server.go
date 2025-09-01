package server

import (
	"context"
	"flow-bridge-mcp/pkg/logger"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
}

func NewMcpServerManager(streamableHttpTransprot *StreamableHttpTransprot, log *logger.Logger) (mcpServerManager *McpServerManager, err error) {
	// 创建 streamable server
	streamableServer, err := server.NewServer(streamableHttpTransprot.StreamableTransport)
	if err != nil {
		return nil, err
	}
	return &McpServerManager{
		Server:                  streamableServer,
		log:                     log,
		streamableHttpTransprot: streamableHttpTransprot,
	}, err
}
func (m *McpServerManager) Run(ctx context.Context) error {
	// 启动 MCP 服务器
	serverErrChan := make(chan error, 1)
	go func() {
		m.log.Info("Starting MCP server")
		serverErrChan <- m.Server.Run()
	}()
	select {
	case err := <-serverErrChan:
		m.log.Error("MCP server error: %v", err)
		return err
	case <-ctx.Done():
		m.log.Info("Shutting down MCP server")
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
