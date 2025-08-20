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
) *gin.Engine {

	// 作为mcp服务对外提供服务
	router := app.app.Group("/")
	{
		router.GET("/sse/:serverId", OpenapiService.Create)
		router.POST("/message", OpenapiService.Create)
		router.POST("/mcp", OpenapiService.Create)
	}
	// 作为api服务对外提供服务
	router = router.Group("/v1")
	{
		//router.GET("/", HomeServer.Index)
		router.POST("/openapi", OpenapiService.Create)
	}

	return app.app
}
