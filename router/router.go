package router

import (
	"flow-bridge-mcp/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewApp, NewRouter)

func NewRouter(
	app *App,
	UserServer *service.UserService,
	HomeServer *service.HomeService,
	LoginService *service.LoginService,
) *gin.Engine {

	// 作为mcp服务对外提供服务
	router := app.app.Group("/")
	{
		router.GET("/sse", HomeServer.Index)
		router.POST("/message", UserServer.Login)
		router.POST("/mcp", LoginService.Login)
	}
	// 作为api服务对外提供服务
	router = router.Group("/v1")
	{
		router.GET("/", HomeServer.Index)
		router.GET("/login", UserServer.Login)
		router.POST("/user/login", LoginService.Login)
	}

	return app.app
}
