package data

import (
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

func (o *OpenapiRepo) Create(serverInfo *model.McpServer) (err error) {

	return
}
