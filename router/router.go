package router

import (
	"flow-bridge-mcp/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewApp, NewRouter)

func NewRouter(
	app *App,
	OpenapiService *service.OpenapiService,
	serverService *service.McpServerService,
	toosService *service.McpToosService,
	mcpGateWayService *service.McpGatewayService,
) *gin.Engine {

	// 作为mcp服务对外提供服务
	router := app.app.Group("/")
	{
		//router.GET("/sse/:serverId", OpenapiService.Upload)
		//router.POST("/message", OpenapiService.Upload)
		router.POST("/gateway/mcp", mcpGateWayService.McpStreamable)
	}
	// 作为api服务对外提供服务
	router = router.Group("/v1")
	{
		router.POST("/openapi/upload", OpenapiService.Upload)
		router.POST("/openapi/updateForAuth", OpenapiService.UpdateForAuth)
		router.POST("/mcpServer/updateByUUID", serverService.UpdateMcpServerByUUID)
		router.POST("/mcpServer/getMcpServerTools", toosService.GetMcpServerTools)
		router.POST("/mcpServer/getMcpConnectTokenByUUID", serverService.GetMcpConnectTokenByUUID)
	}

	return app.app
}
