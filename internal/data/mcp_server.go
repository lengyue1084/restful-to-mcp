package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
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

// GetDb biz 层开启事务可以使用该方法获取同一个tx对象，避免多表操作时，tx对象不一致
func (o *McpServerRepo) GetDb(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return o.data.Db
}

func (o *McpServerRepo) Create(ctx context.Context, serverInfo *model.McpServer) (err error) {
	err = o.data.Db.Create(serverInfo).Error
	if err != nil {
		o.log.ErrorWithContext(ctx, "create mcp server error: %v", err)
		return
	}
	return
}

func (o *McpServerRepo) UpdateByUUID(ctx context.Context, serverInfo *model.McpServer) (err error) {
	return
}
