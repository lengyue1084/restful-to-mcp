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
func (m *McpServerService) GetMcpServerInfoByUUID(c *gin.Context) {
	var req api.GetMcpServerInfoByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "GetMcpServerInfoByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := m.msUc.GetMcpServerInfoByUUID(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "GetMcpServerInfoByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("获取失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "获取成功", resp)

}

func (m *McpServerService) UpdateMcpServerByUUID(c *gin.Context) {
	var req *api.UpdateMcpServerByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "UpdateMcpServerByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := m.msUc.UpdateMcpServerByUUID(c, req.UUID, req.Name, req.Description)
	if err != nil {
		m.log.ErrorWithContext(c, "UpdateMcpServerByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("更新失败，err:%+v", err), nil)
		return
	}
	response.Success(c, "更新成功", resp)
}

func (m *McpServerService) GetMcpConnectTokenByUUID(c *gin.Context) {

	var req api.GetMcpConnectTokenByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := m.msUc.GetMcpConnectTokenByUUID(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "GetMcpConnectTokenByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("获取失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "获取成功", resp)
}

func (m *McpServerService) CreateMcpServerByForm(c *gin.Context) {

	var req *api.CreateMcpServerByFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "CreateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := m.msUc.CreateMcpServerByForm(c, req)
	if err != nil {
		m.log.ErrorWithContext(c, "CreateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("创建失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "创建成功", resp)
}

func (m *McpServerService) DeleteMcpServerByUUID(c *gin.Context) {
	var req *api.DeleteMcpServerByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "DeleteMcpServerByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	err := m.msUc.DeleteMcpServerByUUID(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "DeleteMcpServerByUUID error: %+v", err)
		response.Error(c, fmt.Sprintf("删除失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "删除成功", nil)

}

func (m *McpServerService) UpdateMcpServerByForm(c *gin.Context) {
	var req *api.UpdateMcpServerByFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.log.ErrorWithContext(c, "UpdateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), err)
		return
	}
	resp, err := m.msUc.UpdateMcpServerByForm(c, req)
	if err != nil {
		m.log.ErrorWithContext(c, "UpdateMcpServerByForm error: %+v", err)
		response.Error(c, fmt.Sprintf("更新失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "更新成功", resp)
}
