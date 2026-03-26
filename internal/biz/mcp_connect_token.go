package biz

import (
	"context"
	"restful-to-mcp/internal/data/model"
)

type McpConnectTokenRepo interface {
	Create(ctx context.Context, serverInfo *model.McpConnectToken) (err error)
}
