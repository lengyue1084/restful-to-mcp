package service

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/pkg/response"
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type McpServerService struct {
	msUc *biz.McpServerUseCase
	log  *logger.Logger
}

func NewMcpServerService(msUc *biz.McpServerUseCase, log *logger.Logger) *McpServerService {
	return &McpServerService{
		msUc: msUc,
		log:  log,
	}
}

func (o *McpServerService) UpdateMcpServerByUUID(ctx *gin.Context) {
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

func (o *McpServerService) GetMcpConnectTokenByUUID(c *gin.Context) {

	var req api.GetMcpConnectTokenByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		o.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, "参数错误", err)
		return
	}
	resp, err := o.msUc.GetMcpConnectTokenByUUID(c, req.UUID)
	if err != nil {
		o.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, "获取失败", err)
		return
	}
	response.Success(c, "获取成功", resp)
	return
}
