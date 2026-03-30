package data

import (
	"context"
	"restful-to-mcp/internal/biz"
	"restful-to-mcp/internal/data/database"
	"restful-to-mcp/internal/data/model"
	"restful-to-mcp/pkg/logger"
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

func (m *McpConnectToken) GetByConnectToken(ctx context.Context, connectToken string) (tokenInfo *model.McpConnectToken, err error) {
	tokenInfo = &model.McpConnectToken{}
	err = m.data.Db.WithContext(ctx).
		Where("connect_token = ?", connectToken).
		Order("id DESC").
		Find(tokenInfo).Error
	if err != nil {
		m.log.Error("get mcp connect token error: %v", err)
		return nil, err
	}
	return tokenInfo, nil
}
