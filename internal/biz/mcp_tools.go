package biz

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
)

type McpToolsRepo interface {
	Create(ctx context.Context, mcpToolInfo *model.McpTools) (err error)
	CreateMcpToolsBatch(ctx context.Context, mcpServerId int64, uuid string, allTools []string, mcpToolInfo []*model.McpTools) (err error)
	UpdateToolsForAuthWithTx(ctx context.Context, uuid string, tools []*api.Tools) (err error)
}
