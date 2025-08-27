package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
	"gorm.io/gorm/clause"
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
	err = m.data.Db.WithContext(ctx).Create(mcpToolInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "create mcp tools error: %v", err)
		return
	}
	return
}

func (m *McpToolsRepo) CreateMcpToolsBatch(ctx context.Context, mcpServerId int64, uuid string, allTools []string, mcpToolInfo []*model.McpTools) (err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch get tx error: %v", err)
		return
	}

	// 1、批量删除本次没有包含的工具
	err = db.WithContext(ctx).Where("mcp_server_id = ? and uuid = ? and name not in ?", mcpServerId, uuid, allTools).
		Delete(&model.McpTools{}).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch delete error: %v", err)
		return
	}
	// 2、查询工具是否存在，如果不存在的插入、如果存在的更新
	err = db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "mcp_server_id"},
			{Name: "uuid"},
			{Name: "name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"id",
			"description",
			"mcp_server_type",
			"method",
			"endpoint",
			"headers",
			"args",
			"request_body",
			"response_body",
			"input_schema",
			"annotations",
			"security",
			"is_auth",
			"auth_mode",
			"is_platform_auth",
			"is_show",
		}),
	}).CreateInBatches(mcpToolInfo, 100).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch create error: %v", err)
		return
	}
	return
}
