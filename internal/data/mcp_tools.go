package data

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	_const "flow-bridge-mcp/pkg/const"
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
	err = db.WithContext(ctx).Where("mcp_server_id = ? and mcp_server_uuid = ? and name not in ?", mcpServerId, uuid, allTools).
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
			Where("mcp_server_id = ? AND mcp_server_uuid = ? AND name = ?",
				tool.McpServerId, tool.McpServerUUID, tool.Name).
			Find(&mcpTool).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch find error: %v", err)
			return
		}
		if mcpTool.ID == 0 {
			// 插入
			var toolInfo model.McpTools
			err = db.WithContext(ctx).Where("name = ?", tool.Name).Find(&toolInfo).Error
			if err != nil {
				m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch find serial_number error: %v", err)
				return
			}
			if toolInfo.ID > 0 {
				tool.IsRepeat = _const.CommonStatusYes //重复
			}
			err = db.WithContext(ctx).Create(tool).Error
			if err != nil {
				m.log.ErrorWithContext(ctx, "CreateMcpToolsBatch create error: %v", err)
				return
			}
			continue
		}
		err = db.WithContext(ctx).Model(&model.McpTools{}).
			Where("id = ?", mcpTool.ID).Select("*").
			Omit("SerialNumber", "McpServerId", "McpServerType", "McpServerUUID", "ID", "uuid", "CreatedAt", "DeletedAt").
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
		return
	}
	for _, tool := range tools {
		err = db.WithContext(ctx).Model(&model.McpTools{}).
			Where("id = ? and mcp_server_uuid = ?", tool.ID, uuid).
			Update("is_platform_auth", tool.IsAuth).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateToolsForAuthWithTx update error: %v", err)
			return
		}
	}
	return
}

func (m *McpToolsRepo) GetMcpServerToolsByServerUUID(ctx context.Context, uuid string) (mcpTools []*model.McpTools, err error) {
	err = m.data.Db.WithContext(ctx).Where("mcp_server_uuid = ?", uuid).Find(&mcpTools).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server tools by uuid error: %v", err)
		return
	}
	return
}

func (m *McpToolsRepo) CreateMcpServerTool(ctx context.Context, mcpToolInfo *model.McpTools) (uuid string, err error) {
	err = m.data.Db.WithContext(ctx).Create(mcpToolInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "create mcp server tool error: %v", err)
		return
	}
	uuid = mcpToolInfo.UUID
	return
}

func (m *McpToolsRepo) GetMcpServerToolByNameWithUUID(ctx context.Context, uuid string, name string) (tool *model.McpTools, err error) {

	err = m.data.Db.WithContext(ctx).Where("mcp_server_uuid = ? and name = ?", uuid, name).Find(&tool).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server tools by name with uuid error: %v", err)
		return
	}
	return
}

func (m *McpToolsRepo) GetMcpServerToolByName(ctx context.Context, name string) (tool *model.McpTools, err error) {
	err = m.data.Db.WithContext(ctx).Where("name = ?", name).Find(&tool).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server tools by name error: %v", err)
		return
	}
	return
}

func (m *McpToolsRepo) GetMcpServerToolInfoByUUID(ctx context.Context, uuid string) (tool *model.McpTools, err error) {
	err = m.data.Db.WithContext(ctx).Where("uuid = ?", uuid).Find(&tool).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server tools info by uuid error: %v", err)
		return
	}
	return

}

func (m *McpToolsRepo) UpdateMcpServerTool(ctx context.Context, tool *model.McpTools, uuid string) (err error) {
	err = m.data.Db.WithContext(ctx).Model(&model.McpTools{}).
		Where("uuid = ?", uuid).
		Select([]string{
			"Name",
			"Description",
			"Method",
			"Endpoint",
			"Headers",
			"Args",
			"ToolSchema",
			"Annotations",
			"security",
			"IsPlatformAuth",
			"IsAuth",
			"AuthMode",
			"IsRepeat",
		}).
		Updates(tool).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "update mcp server tool error: %v", err)
		return
	}

	return
}
