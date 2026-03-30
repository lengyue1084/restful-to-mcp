package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"restful-to-mcp/api"
	"restful-to-mcp/internal/biz"
	mcpServer "restful-to-mcp/internal/mcp/server"
	"restful-to-mcp/internal/pkg/response"
	"restful-to-mcp/pkg/logger"
)

type McpToosService struct {
	mtUc             *biz.McpToolsUserCase
	oaUc             *biz.OpenapiUseCase
	mcpServerManager *mcpServer.McpServerManager
	log              *logger.Logger
}

func NewMcpToosService(mtUc *biz.McpToolsUserCase, oaUc *biz.OpenapiUseCase, mcpServerManager *mcpServer.McpServerManager, log *logger.Logger) *McpToosService {
	return &McpToosService{
		mtUc:             mtUc,
		oaUc:             oaUc,
		mcpServerManager: mcpServerManager,
		log:              log,
	}
}

func (m *McpToosService) refreshTools(c *gin.Context) {
	m.oaUc.UpdateToolsForCache(c)
	m.mcpServerManager.RegisterToolFromCache()
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
}

func (m *McpToosService) GetMcpServerToolsByUUID(c *gin.Context) {
	var req *api.GetMcpServerToolsByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	resp, err := m.mtUc.GetMcpServerToolsByUUID(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "查询工具列表失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("查询工具列表失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "查询成功", resp)
}

func (m *McpToosService) CreateMcpServerTool(c *gin.Context) {
	var req *api.CreateMcpServerToolRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	resp, err := m.mtUc.CreateMcpServerTool(c, req)
	if err != nil {
		m.log.ErrorWithContext(c, "创建工具失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("创建工具失败,err:%+v", err), nil)
		return
	}
	m.refreshTools(c)
	response.Success(c, "创建成功", resp)
}

func (m *McpToosService) UpdateMcpServerTool(c *gin.Context) {
	var req *api.UpdateMcpServerToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	if err := api.ValidMethods(req.Method); err != nil {
		response.Error(c, fmt.Sprintf("method参数错误,err:%+v", err), nil)
		return
	}

	m.oaUc.UpdateToolsForOldCache(c)
	resp, err := m.mtUc.UpdateMcpServerTool(c, req)
	if err != nil {
		m.log.ErrorWithContext(c, "更新工具失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("更新工具失败,err:%+v", err), nil)
		return
	}
	m.refreshTools(c)
	response.Success(c, "更新成功", resp)
}

func (m *McpToosService) GetToolsInfoByUUID(c *gin.Context) {
	var req *api.GetToolsInfoByUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	resp, err := m.mtUc.GetToolsInfoByUUID(c, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(c, "查询工具失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("查询工具失败,err:%+v", err), nil)
		return
	}
	response.Success(c, "查询成功", resp)
}

func (m *McpToosService) TestMcpServerTool(c *gin.Context) {
	var req *api.TestMcpServerToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}

	resp, err := m.mtUc.TestMcpServerTool(c, req)
	if err != nil {
		m.log.ErrorWithContext(c, "测试工具失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("测试工具失败,err:%+v", err), nil)
		return
	}

	response.Success(c, "测试成功", resp)
}
