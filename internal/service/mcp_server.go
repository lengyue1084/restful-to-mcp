package service

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/pkg/response"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
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
		response.Error(ctx, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := o.msUc.UpdateMcpServerByUUID(ctx, req.UUID, req.Name, req.Description)
	if err != nil {
		o.log.ErrorWithContext(ctx, "UpdateMcpServerByUUID error: %+v", err)
		response.Error(ctx, fmt.Sprintf("更新失败，err:%+v", err), nil)
		return
	}
	response.Success(ctx, "更新成功", resp)

}

func (o *McpServerService) GetMcpConnectTokenByUUID(c *gin.Context) {

	var req api.GetMcpConnectTokenByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		o.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := o.msUc.GetMcpConnectTokenByUUID(c, req.UUID)
	if err != nil {
		o.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("获取失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "获取成功", resp)
	return
}

func (o *McpServerService) CreateMcpServerByForm(c *gin.Context) {

	var req *api.CreateMcpServerByFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		o.log.ErrorWithContext(c, "CreateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := o.msUc.CreateMcpServerByForm(c, req)
	if err != nil {
		o.log.ErrorWithContext(c, "CreateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("创建失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "创建成功", resp)
	return
}
