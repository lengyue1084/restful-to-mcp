package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
)

var _ biz.McpToolsRepo = (*McpToolsRepo)(nil)

type McpToolsRepo struct {
	data *database.Data
	log  *logger.Logger
}

func NewMcpToolsRepo(data *database.Data, log *logger.Logger) biz.McpToolsRepo {
	return &McpToolsRepo{
		data: data,
		log:  log,
	}
}

func (m *McpToolsRepo) Create(ctx context.Context, mcpToolInfo *model.McpTools) (err error) {

	return
}
