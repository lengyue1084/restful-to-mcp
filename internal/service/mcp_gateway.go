package service

import (
	"context"
	"encoding/json"
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

	var inputStr = `{
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"description": "Name of pet that needs to be updated"
		},
		"petId": {
			"type": "integer",
			"description": "ID of pet that needs to be updated"
		},
		"status": {
			"type": "string",
			"description": "Status of pet that needs to be updated"
		}
	},
	"required": ["petId"]
}`
	var toolSchema protocol.InputSchema
	err := json.Unmarshal([]byte(inputStr), &toolSchema)
	if err != nil {
		m.log.Errorf("Failed to unmarshal tool schema: %v", err)
		c.JSON(500, gin.H{"error": "Failed to unmarshal tool schema"})
		return
	}

	toolInfo := &protocol.Tool{
		Name:           "createUsersWithListInput",
		Description:    "Creates list of users with given input array.",
		InputSchema:    toolSchema,
		OutputSchema:   protocol.OutputSchema{},
		Annotations:    nil,
		RawInputSchema: nil,
	}
	m.mcpServerManager.Server.RegisterTool(toolInfo, handleTimeRequest)
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
