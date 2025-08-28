package service

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/pkg/response"
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type McpServerSverService struct {
	msUc *biz.McpServerUseCase
	log  *logger.Logger
}

func NewMcpServerService(msUc *biz.McpServerUseCase, log *logger.Logger) *McpServerSverService {
	return &McpServerSverService{
		msUc: msUc,
		log:  log,
	}
}

func (o *McpServerSverService) UpdateMcpServerByUUID(ctx *gin.Context) {
	var req *api.UpdateMcpServerByUUIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		o.log.ErrorWithContext(ctx, "UpdateMcpServerByUUID error: %+v", err)
		response.Error(ctx, "参数错误", err)
		return
	}
	resp, err := o.msUc.UpdateMcpServerByUUID(ctx, req.UUID, req.Name, req.Description)
	if err != nil {
		o.log.ErrorWithContext(ctx, "UpdateMcpServerByUUID error: %+v", err)
		response.Error(ctx, "更新失败", nil)
		return
	}
	response.Success(ctx, "更新成功", resp)

}
