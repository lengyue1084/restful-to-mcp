package biz

import (
	"context"
	"flow-bridge-mcp/internal/data/model"
)

type McpToolsRepo interface {
	Create(ctx context.Context, mcpToolInfo *model.McpTools) (err error)
}
