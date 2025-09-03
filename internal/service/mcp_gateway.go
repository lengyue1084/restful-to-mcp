package service

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

type McpGatewayService struct {
	McpGateWayUc     *biz.McpGatewayUseCase
	log              *logger.Logger
	mcpServerManager *mcpServer.McpServerManager
	mcpServerUseCase *biz.McpServerUseCase
}

func NewMcpGatewayService(msUc *biz.McpGatewayUseCase, mcpServerManager *mcpServer.McpServerManager, log *logger.Logger) *McpGatewayService {
	// 预启动服务器
	ctx := context.Background()
	//defer cancel()
	go func() {
		if err := mcpServerManager.Run(ctx); err != nil {
			log.Errorf("MCP server failed: %v", err)
		}
	}()
	return &McpGatewayService{
		McpGateWayUc:     msUc,
		log:              log,
		mcpServerManager: mcpServerManager,
	}
}

func (m *McpGatewayService) McpStreamable(c *gin.Context) {
	platformToken := c.GetHeader("platform-token")
	serviceToken := c.GetHeader("service-token")
	fmt.Printf("platformToken:%s\n,serviceToken:%s\n", platformToken, serviceToken)

	// 使用已预启动的服务器管理器处理连接
	m.mcpServerManager.HandleConnection(c)

}
