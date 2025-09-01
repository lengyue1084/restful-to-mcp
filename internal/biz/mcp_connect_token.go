package biz

import (
	"context"
	"flow-bridge-mcp/internal/data/model"
)

type McpConnectTokenRepo interface {
	Create(ctx context.Context, serverInfo *model.McpConnectToken) (err error)
}
