package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
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
	err = db.WithContext(ctx).Where("uuid = ? and status = ?", UUID, model.StatusHidden).Find(&serverInfo).Error
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
	err = db.WithContext(ctx).Where("uuid = ? and status = ?", serverInfo.UUID, model.StatusHidden).Find(&serverMcpInfo).Error
	if err != nil {
		m.log.ErrorWithContext(ctx, "get mcp server error: %v", err)
		return
	}
	var serverMcpInfoId = serverMcpInfo.ID
	if serverMcpInfo.ID > 0 {
		//更新
		err = db.WithContext(ctx).
			Where("uuid = ? and id = ?", serverInfo.UUID, serverMcpInfo.ID).
			Omit("ServiceToken", "PlatformToken").
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

func (m *McpServerRepo) UpdateByUUID(ctx context.Context, serverInfo *model.McpServer) (err error) {
	return
}
