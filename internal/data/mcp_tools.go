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

func (m *McpToolsRepo) CreateMcpToolsBatch(ctx context.Context, mcpServerId int64, uuid string, allTools []string, mcpToolInfo []*model.McpTools) (err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch get tx error: %v", err)
		return
	}

	// 1、批量删除本次没有包含的工具
	err = db.Where("mcp_server_id = ? and uuid = ? and name in (?)", mcpServerId, uuid, allTools).
		Delete(&model.McpTools{}).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch delete error: %v", err)
		return
	}

	// 2、本次包含的工具，更新操作
	// 3、本次没有包含的工具，插入操作
	err = db.WithContext(ctx).Create(mcpToolInfo).Error

	return
}
