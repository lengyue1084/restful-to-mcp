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
	sverService *service.McpServerSverService,
) *gin.Engine {

	// 作为mcp服务对外提供服务
	router := app.app.Group("/")
	{
		router.GET("/sse/:serverId", OpenapiService.Upload)
		router.POST("/message", OpenapiService.Upload)
		router.POST("/mcp/serverId", OpenapiService.Upload)
	}
	// 作为api服务对外提供服务
	router = router.Group("/v1")
	{
		router.POST("/openapi/upload", OpenapiService.Upload)
		router.POST("/openapi/updateForAuth", OpenapiService.UpdateForAuth)
		router.POST("/mcpServer/updateByUUID", sverService.UpdateMcpServerByUUID)
	}

	return app.app
}
