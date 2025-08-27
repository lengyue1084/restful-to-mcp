package biz

import (
	"context"
	"flow-bridge-mcp/internal/data/model"
)

type McpServerRepo interface {
	Create(ctx context.Context, serverInfo *model.McpServer) (err error)
	UpdateByUUID(ctx context.Context, serverInfo *model.McpServer) (err error)
	GetMcpServerInfo(ctx context.Context, UUID string) (serverInfo *model.McpServer, err error)
	CreateWithTx(ctx context.Context, serverInfo *model.McpServer) (id int64, err error)
}
