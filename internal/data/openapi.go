package data

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
)

var _ biz.OpenapiRepo = (*OpenapiRepo)(nil)

type OpenapiRepo struct {
	data *Data
	log  *logger.Logger
}

func NewOpenapiRepo(d *Data, log *logger.Logger) biz.OpenapiRepo {
	return &OpenapiRepo{
		data: d,
		log:  log,
	}
}

func (o *OpenapiRepo) Create(ctx context.Context, serverInfo *model.McpServer) (err error) {
	err = o.data.db.Where("namesss = ?", serverInfo.CreatedAt).Updates(serverInfo).Error
	if err != nil {
		o.log.ErrorWithContext(ctx, "create mcp server error: %v", err)
		return
	}
	return
}
