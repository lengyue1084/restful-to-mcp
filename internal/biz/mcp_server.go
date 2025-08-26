package biz

import (
	"context"
	"flow-bridge-mcp/internal/data/model"
)

type McpServerRepo interface {
	Create(ctx context.Context, serverInfo *model.McpServer) (err error)
	UpdateByUUID(ctx context.Context, serverInfo *model.McpServer) (err error)
}
