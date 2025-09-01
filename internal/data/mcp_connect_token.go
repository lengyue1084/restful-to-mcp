package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
)

type McpConnectToken struct {
	data *database.Data
	log  *logger.Logger
}

func NewMcpConnectToken(data *database.Data, log *logger.Logger) biz.McpConnectTokenRepo {
	return &McpConnectToken{
		data: data,
		log:  log,
	}
}

func (m *McpConnectToken) Create(ctx context.Context, mcpConnectToken *model.McpConnectToken) (err error) {
	err = m.data.Db.WithContext(ctx).Create(mcpConnectToken).Error
	if err != nil {
		m.log.Error("create mcp connect token error: %v", err)
		return
	}
	return
}
