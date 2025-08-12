package router

import (
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	app *gin.Engine
}

func NewApp(
	middleware *middleware.Middleware,
	conf *conf.Conf,
	log *conf.Logger,
) *App {
	if !conf.Conf.GetBool("server.dev") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	//	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	//}
	r.Use(middleware.Cors(), middleware.TraceId(),
		middleware.ZapLogger(log), middleware.Recovery())
	return &App{
		app: r,
	}
}
