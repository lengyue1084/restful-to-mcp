package biz

import (
	"encoding/base64"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	openapi2 "flow-bridge-mcp/internal/mcp/transformer/openapi"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"github.com/gin-gonic/gin"
)

type OpenapiRepo interface {
	Create(serverInfo *model.McpServer) (err error)
}
type OpenapiUseCase struct {
	OpenapiRepo OpenapiRepo
	log         *logger.Logger
}

func NewOpenapiUserCase(repo OpenapiRepo, log *logger.Logger) *OpenapiUseCase {
	return &OpenapiUseCase{
		OpenapiRepo: repo,
		log:         log,
	}
}

func (o *OpenapiUseCase) Create(ctx *gin.Context, req *api.ServerInfoRequest) (resp *api.ServerInfoResponse, err error) {

	decodeString, err := base64.URLEncoding.DecodeString(req.FileContent)
	if err != nil {
		o.log.ErrorWithContext(ctx, "base64解析失败，err:%+v", err)
		return nil, err
	}
	req.FileContent = string(decodeString)
	converter := openapi2.NewConverter()
	mcpConfig, err := converter.Convert(decodeString)
	if err != nil {
		o.log.ErrorWithContext(ctx, "数据转换错误，err:%+v", err)
		return nil, err
	}
	mcpConfig.Tenant = tool.WithPrefix(mcpConfig.Tenant, req.Name)

	var serverInfo *model.McpServer
	if err := tool.Copy(&serverInfo, req); err != nil {
		o.log.ErrorWithContext(ctx, "数据转换错误，err:%+v", err)
		return nil, err
	}

	err = o.OpenapiRepo.Create(serverInfo)
	return
}
