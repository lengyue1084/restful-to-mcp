package biz

import (
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
)

type OpenapiRepo interface {
	Create(serverInfo *model.McpServer) (err error)
}
type OpenapiUseCase struct {
	OpenapiRepo OpenapiRepo
	log         *logger.Logger
}

func NewOpenapiUserCase(repo OpenapiRepo, log *logger.Logger) *OpenapiUseCase {
	return &OpenapiUseCase{
		OpenapiRepo: repo,
		log:         log,
	}
}
