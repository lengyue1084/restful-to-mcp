package service

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type McpGatewayService struct {
	McpGateWayUc     *biz.McpGatewayUseCase
	log              *logger.Logger
	mcpServerManager *mcpServer.McpServerManager
	mcpServerUseCase *biz.McpServerUseCase
	cache            *cache.MemoryCache
}

func NewMcpGatewayService(msUc *biz.McpGatewayUseCase, mcpServerUseCase *biz.McpServerUseCase, mcpServerManager *mcpServer.McpServerManager, log *logger.Logger, cache *cache.MemoryCache) *McpGatewayService {
	// 预启动服务器
	ctx := context.Background()
	go func() {
		if err := mcpServerManager.Run(ctx); err != nil {
			log.Errorf("MCP server failed: %v", err)
		}
	}()
	return &McpGatewayService{
		McpGateWayUc:     msUc,
		log:              log,
		mcpServerManager: mcpServerManager,
		cache:            cache,
		mcpServerUseCase: mcpServerUseCase,
	}
}

func (m *McpGatewayService) McpStreamable(c *gin.Context) {

	ctx := c.Request.Context()
	// 只处理需要的 header
	if platformToken := c.GetHeader(_const.PlatformToken); platformToken != "" {
		ctx = context.WithValue(ctx, _const.PlatformToken, platformToken)
	}
	if serviceToken := c.GetHeader(_const.ServiceToken); serviceToken != "" {
		ctx = context.WithValue(ctx, _const.ServiceToken, serviceToken)
	}
	if traceId := c.Value(_const.TraceId); traceId != "" {
		ctx = context.WithValue(ctx, _const.TraceId, traceId)
	}
	if spanId := c.Value(_const.SpanId); spanId != "" {
		ctx = context.WithValue(ctx, _const.SpanId, spanId)
	}

	serverToken := c.Param(_const.ServerToken)
	if serverToken != "" {
		ctx = context.WithValue(ctx, _const.ServerToken, serverToken)
	}

	c.Request = c.Request.WithContext(ctx)
	// 使用已预启动的服务器管理器处理连接
	m.mcpServerManager.HandleConnection(c)

}
