package data

import (
	"context"
	"flow-bridge-mcp/api"
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
	// 然后单独更新已存在的记录
	for _, tool := range mcpToolInfo {
		var mcpTool model.McpTools
		err = db.WithContext(ctx).Model(&model.McpTools{}).
			Where("mcp_server_id = ? AND uuid = ? AND name = ?",
				tool.McpServerId, tool.UUID, tool.Name).
			Find(&mcpTool).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch find error: %v", err)
			return
		}
		if mcpTool.ID == 0 {
			// 插入
			err = db.WithContext(ctx).Create(tool).Error
			if err != nil {
				m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch create error: %v", err)
				return
			}
			continue
		}
		err = db.WithContext(ctx).Model(&model.McpTools{}).
			Where("id = ?", mcpTool.ID).
			Updates(tool).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch update error: %v", err)
			return
		}
	}
	return
}

func (m *McpToolsRepo) UpdateToolsForAuthWithTx(ctx context.Context, uuid string, tools []*api.Tools) (err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateToolsForAuthWithTx get tx error: %v", err)
	}
	for _, tool := range tools {
		err = db.WithContext(ctx).Model(&model.McpTools{}).
			Where("id = ?", tool.ID).
			Update("is_platform_auth", tool.IsAuth).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateToolsForAuthWithTx update error: %v", err)
			return
		}
	}
	return
}

func (m *McpToolsRepo) GetMcpServerTools(ctx context.Context, uuid string) (mcpTools []*model.McpTools, err error) {
	err = m.data.Db.WithContext(ctx).Where("uuid = ?", uuid).Find(&mcpTools).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server tools error: %v", err)
		return
	}
	return
}
