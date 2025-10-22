package router

import (
	"flow-bridge-mcp/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewApp, NewRouter)

func NewRouter(
	app *App,
	openapiService *service.OpenapiService,
	mcpServerService *service.McpServerService,
	mcpToolsService *service.McpToosService,
	mcpGateWayService *service.McpGatewayService,
) *gin.Engine {

	// MCP网关服务路由
	router := app.app.Group("/")
	{
		//router.GET("/sse/:serverId", OpenapiService.Upload)
		//router.POST("/message", OpenapiService.Upload)
		router.POST("/gateway/mcp", mcpGateWayService.McpStreamable)
		router.POST("/gateway/:serverToken/mcp", mcpGateWayService.McpStreamable)
	}
	// API服务路由 (v1版本)
	apiV1 := router.Group("/v1")
	{
		// 基于openapi文档创建mcpServer
		// OpenAPI文档相关
		apiV1.POST("/openapi/upload", openapiService.Upload)
		apiV1.POST("/openapi/updateForAuth", openapiService.UpdateForAuth)

		// MCP Server管理相关
		apiV1.POST("/mcpServer/updateByUUID", mcpServerService.UpdateMcpServerByUUID)
		// 获取MCP Server的工具 符合mcp tools/list的返回结构,也可以可以通过连接mcp之后获取tools/list
		apiV1.POST("/mcpServer/getMcpConnectTokenByUUID", mcpServerService.GetMcpConnectTokenByUUID)
		apiV1.POST("/mcpServer/deleteMcpServerByUUID", mcpServerService.DeleteMcpServerByUUID)
		apiV1.POST("/mcpServer/getMcpServerInfoByUUID", mcpServerService.GetMcpServerInfoByUUID)
		apiV1.POST("/mcpServer/getMcpServerTools", mcpToolsService.GetMcpServerTools)

		// 表单创建MCP Server
		apiV1.POST("/mcpServer/createByForm", mcpServerService.CreateMcpServerByForm)
		apiV1.POST("/mcpServer/updateMcpServerByForm", mcpServerService.UpdateMcpServerByForm)
		apiV1.POST("/mcpServer/getMcpServerToolsByUUID", mcpToolsService.GetMcpServerToolsByUUID)
		apiV1.POST("/mcpServer/createMcpServerTool", mcpToolsService.CreateMcpServerTool)
		apiV1.POST("/mcpServer/updateMcpServerTool", mcpToolsService.UpdateMcpServerTool)

	}
	app.app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 not found"})
	})

	return app.app
}
