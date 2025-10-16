package router

import (
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/pkg/startup"
	"flow-bridge-mcp/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	app *gin.Engine
}

func NewApp(
	middleware *middleware.Middleware,
	conf *conf.Conf,
	initializer *startup.Initializer,
) *App {
	if !conf.Conf.GetBool("server.dev") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	timeoutSecond := conf.Conf.GetDuration("server.http.timeout")
	if timeoutSecond <= 0 {
		timeoutSecond = conf.Conf.GetDuration("server.timeout")
	}

	r.Use(
		middleware.TraceIdMiddleware(),              // 追踪ID - 最优先
		middleware.CorsMiddleware(),                 // CORS - 尽早处理，避免不必要的处理
		middleware.RecoveryMiddleware(),             // panic恢复 - 尽早放置，捕获所有panic
		middleware.TimeoutMiddleware(timeoutSecond), // 超时控制 - 在主要业务逻辑前
		middleware.LoggerMiddleware(),               // 日志记录 - 记录完整处理过程
	)
	err := initializer.Initialize()
	if err != nil {
		initializer.Log.Infof("更新全量mcp工具数据失败: %v", err)
	}
	initializer.Log.Info("加载全量mcp工具成功")
	return &App{
		app: r,
	}
}
