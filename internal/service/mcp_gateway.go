package service

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gin-gonic/gin"
	"time"
)

type McpGatewayService struct {
	McpGateWayUc     *biz.McpGatewayUseCase
	log              *logger.Logger
	mcpServerManager *mcpServer.McpServerManager
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

	type currentTimeInputSchema struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Data        struct {
			Time int64 `json:"time"`
		} `json:"object"`
	}

	// 注册示例工具
	tool, err := protocol.NewTool("current_time", "获取指定时区的当前时间", &currentTimeInputSchema{})
	if err != nil {
		m.log.Errorf("Failed to create tool: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create tool"})
		return
	}
	m.mcpServerManager.Server.RegisterTool(tool, handleTimeRequest)
	//m.mcpServerManager.Server.UnregisterTool("current_time")
	// 使用已预启动的服务器管理器处理连接
	m.mcpServerManager.HandleConnection(c)

}

func handleTimeRequest(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("无效的时区: %v", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: time.Now().In(loc).String(),
			},
		},
	}, nil
}
