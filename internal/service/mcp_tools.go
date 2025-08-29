package service

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/pkg/response"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

type McpToosService struct {
	mtUc *biz.McpToolsUserCase
	log  *logger.Logger
}

func NewMcpToosService(mtUc *biz.McpToolsUserCase, log *logger.Logger) *McpToosService {
	return &McpToosService{
		mtUc: mtUc,
		log:  log,
	}
}

func (m *McpToosService) GetMcpServerTools(c *gin.Context) {

	var req *api.GetMcpServerToolsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	resp, err := m.mtUc.GetMcpServerTools(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "查询工具列表失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("查询工具列表失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "查询成功", resp.Tools)
	return

}
