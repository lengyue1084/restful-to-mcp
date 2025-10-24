package data

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"gorm.io/gorm"
)

var _ biz.McpServerRepo = (*McpServerRepo)(nil)

type McpServerRepo struct {
	data *database.Data
	log  *logger.Logger
}

func NewMcpServerRepo(d *database.Data, log *logger.Logger) biz.McpServerRepo {
	return &McpServerRepo{
		data: d,
		log:  log,
	}
}

func (m *McpServerRepo) Create(ctx context.Context, serverInfo *model.McpServer) (err error) {
	err = m.data.Db.WithContext(ctx).Create(serverInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "create mcp server error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) GetMcpServerInfo(ctx context.Context, UUID string) (serverInfo *model.McpServer, err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "get tx error: %v", err)
		return
	}
	err = db.WithContext(ctx).Where("uuid = ? and status = ?", UUID, _const.StatusHidden).Find(&serverInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) CreateWithTx(ctx context.Context, serverInfo *model.McpServer) (id int64, err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "get tx error: %v", err)
		return
	}
	var serverMcpInfo model.McpServer
	err = db.WithContext(ctx).Where("uuid = ?", serverInfo.UUID).Find(&serverMcpInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	var serverMcpInfoId = serverMcpInfo.ID
	if serverMcpInfo.ID > 0 {
		//更新
		err = db.WithContext(ctx).
			Where("uuid = ? and id = ?", serverInfo.UUID, serverMcpInfo.ID).
			Select("*").
			Omit("ServiceToken", "PlatformToken", "SerialNumber", "ID", "UUID", "CreatedAt", "source").
			Updates(serverInfo).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "update mcp server error: %v", err)
			return
		}

	} else {
		//新增
		err = db.WithContext(ctx).Create(&serverInfo).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "create mcp server error: %v", err)
			return
		}
		serverMcpInfoId = serverInfo.ID
	}
	return int64(serverMcpInfoId), nil
}

func (m *McpServerRepo) GetMcpServerInfoByID(ctx context.Context, id int64) (mcpServerInfo *model.McpServer, err error) {

	err = m.data.Db.WithContext(ctx).Where("id = ?", id).
		Preload("Tools").
		Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) GetMcpServerInfoByUUID(ctx context.Context, id string) (mcpServerInfo *model.McpServer, err error) {

	err = m.data.Db.WithContext(ctx).Where("uuid = ?", id).
		Preload("Tools").
		Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) UpdateMcpServerForAuthWithTx(ctx context.Context, uuid string, isAuth _const.AuthTypeStatus, serviceToken, platformToken string) (err error) {
	db, err := m.data.GetDb(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "get tx error: %v", err)
		return
	}
	var mcpServerInfo = &model.McpServer{}
	err = db.WithContext(ctx).Where("uuid = ?", uuid).Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	if mcpServerInfo.ID <= 0 {
		return fmt.Errorf("没有查询到该server信息")
	}
	err = db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Select("ServiceToken", "PlatformToken", "IsAuth", "Status").
		Updates(model.McpServer{
			ServiceToken:  serviceToken,
			PlatformToken: platformToken,
			IsAuth:        isAuth,
			Status:        _const.ServerHadSetToken,
		}).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "update mcp server error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) UpdateMcpServerByUUID(ctx context.Context, uuid, name, description string) (resp *api.UpdateMcpServerByUUIDResponse, err error) {
	var mcpServerInfo = &model.McpServer{}
	err = m.data.Db.WithContext(ctx).Where("uuid = ?", uuid).Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	if mcpServerInfo.ID <= 0 {
		return nil, fmt.Errorf("没有查询到该server信息")
	}

	updateMap := make(map[string]interface{})
	if name != "" {
		updateMap["name"] = name
	}
	if description != "" {
		updateMap["description"] = description
	}
	err = m.data.Db.WithContext(ctx).
		Model(&model.McpServer{}).
		Where("uuid = ?", uuid).
		Updates(updateMap).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "update mcp server error: %v", err)
		return
	}
	return

}

func (m *McpServerRepo) DeleteMcpServerByUUID(ctx context.Context, uuid string) (err error) {
	err = m.data.Db.Transaction(func(tx *gorm.DB) error {
		var mcpServerInfo model.McpServer
		err = tx.Where("uuid = ?", uuid).Find(&mcpServerInfo).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "get mcp server by uuid error: %v", err)
			return err
		}
		if mcpServerInfo.ID <= 0 {
			return fmt.Errorf("没有查询到该server信息")
		}
		err = tx.WithContext(ctx).
			Where("uuid = ?", uuid).
			Delete(&model.McpServer{}).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "delete mcp server error: %v", err)
			return err
		}
		err = tx.WithContext(ctx).
			Where("mcp_server_id = ?", mcpServerInfo.ID).
			Delete(&model.McpTools{}).Error
		if err != nil {
			m.log.ErrorWithContext(ctx, "delete mcp tools error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		m.log.ErrorWithContext(ctx, "delete mcp server error: %v", err)
		return err
	}

	return
}

func (m *McpServerRepo) GetCountMcpServerInfoBySerialNumber(ctx context.Context, serialNumber string) (count int64, err error) {
	err = m.data.Db.WithContext(ctx).
		Model(&model.McpServer{}).
		Where("serial_number = ?", serialNumber).
		Count(&count).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetCountMcpServerInfoBySerialNumber error: %v", err)
		return
	}
	return
}

func (m *McpServerRepo) GetMcpServerInfoWithAllTools(ctx context.Context) (serverInfo []*model.McpServer, err error) {
	err = m.data.Db.WithContext(ctx).
		Where("status in ?", []int{int(_const.ServerHadSetToken), int(_const.ServerTokenIsWorking)}).
		Preload("Tools").
		Find(&serverInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoWithTools error: %v", err)
		return
	}

	return
}

func (m *McpServerRepo) CreateMcpServerByForm(ctx context.Context, serverInfo *model.McpServer) (mcpServer *model.McpServer, err error) {

	var mcpServerInfo = &model.McpServer{}
	err = m.data.Db.WithContext(ctx).Where("uuid = ?", serverInfo.UUID).Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	if mcpServerInfo.ID > 0 {
		return nil, fmt.Errorf("该uuid已存在")
	}
	err = m.data.Db.WithContext(ctx).Select("*").Create(serverInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "create mcp server error: %v", err)
		return
	}
	return serverInfo, nil
}

func (m *McpServerRepo) UpdateMcpServerByForm(ctx context.Context, serverInfo *model.McpServer) (mcpServer *model.McpServer, err error) {

	var mcpServerInfo = &model.McpServer{}
	err = m.data.Db.WithContext(ctx).Where("uuid = ?", serverInfo.UUID).Find(&mcpServerInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	if mcpServerInfo.ID <= 0 {
		return nil, fmt.Errorf("没有查询到该server信息")
	}
	err = m.data.Db.Model(&model.McpServer{}).
		Where("uuid = ?", serverInfo.UUID).
		Select("name", "description", "urls", "version", "is_auth", "platform_token", "service_token", "header", "all_tools", "mcp_server_type", "have_tools", "security").
		Updates(serverInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "update mcp server error: %v", err)
		return nil, err
	}
	err = m.data.Db.Where("uuid = ?", serverInfo.UUID).First(&mcpServer).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	return mcpServer, nil
}
