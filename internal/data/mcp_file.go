package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
)

var _ biz.McpFileRepo = (*McpFileRepo)(nil)

type McpFileRepo struct {
	data *database.Data
	log  *logger.Logger
}

func NewMcpFileRepo(data *database.Data, log *logger.Logger) biz.McpFileRepo {
	return &McpFileRepo{
		data: data,
		log:  log,
	}
}

func (m *McpFileRepo) Create(ctx context.Context, serverInfo *model.McpFile) (err error) {
	err = m.data.Db.WithContext(ctx).Create(serverInfo).Error
	if err != nil {
		m.log.Errorf("create mcp file error: %v", err)
		return err
	}
	return
}

func (m *McpFileRepo) GetMcpFileInfoByMd5(ctx context.Context, md5 string) (mcpFileInfo *model.McpFile, err error) {
	err = m.data.Db.WithContext(ctx).Where("md5 = ?", md5).Find(&mcpFileInfo).Error
	if err != nil {
		m.log.Errorf("get mcp file info by md5 error: %v", err)
		return nil, err
	}
	return
}

func (m *McpFileRepo) UpdateMcpFileById(ctx context.Context, serverInfo *model.McpFile) (err error) {
	err = m.data.Db.WithContext(ctx).Where("id = ?", serverInfo.ID).Updates(serverInfo).Error
	if err != nil {
		m.log.Errorf("update mcp file error: %v", err)
		return err
	}
	return

}
