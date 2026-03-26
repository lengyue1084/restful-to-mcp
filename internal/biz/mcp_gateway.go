package biz

import "restful-to-mcp/pkg/logger"

type McpGatewayRepo interface {
}

type McpGatewayUseCase struct {
	log *logger.Logger
}

func NewMcpGatewayUseCase(log *logger.Logger) *McpGatewayUseCase {
	return &McpGatewayUseCase{
		log: log,
	}
}
