package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"restful-to-mcp/internal/biz"
	mcpServer "restful-to-mcp/internal/mcp/server"
	"restful-to-mcp/internal/pkg/cache"
	_const "restful-to-mcp/pkg/const"
	"restful-to-mcp/pkg/logger"
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
	// 添加请求头信息
	ctx = context.WithValue(ctx, _const.PlatformToken, c.GetHeader(_const.PlatformToken))
	ctx = context.WithValue(ctx, _const.ServiceToken, c.GetHeader(_const.ServiceToken))

	if traceId := c.Value(_const.TraceId); traceId != "" {
		ctx = context.WithValue(ctx, _const.TraceId, traceId)
	}
	if spanId := c.Value(_const.SpanId); spanId != "" {
		ctx = context.WithValue(ctx, _const.SpanId, spanId)
	}

	serverToken := c.Param(_const.ServerPathToken)
	if serverToken != "" {
		resolvedServerUUID, err := m.mcpServerUseCase.ResolveServerToken(ctx, serverToken)
		if err != nil {
			m.log.ErrorWithContext(c, "resolve server token error: %v", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "invalid server token",
			})
			return
		}
		ctx = context.WithValue(ctx, _const.ServerPathToken, resolvedServerUUID)
	}

	c.Request = c.Request.WithContext(ctx)
	// 使用已预启动的服务器管理器处理连接
	m.mcpServerManager.HandleConnection(c)

}
