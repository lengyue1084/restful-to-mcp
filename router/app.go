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
	//gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	//	log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	//}
	// 设置可信代理
	//if conf.Conf.GetBool("server.dev") {
	//	// 开发环境信任本地代理
	//	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	//} else {
	//	// 生产环境根据配置设置
	//	// 如果配置文件中有 trusted_proxies，则使用配置的值
	//	// 否则禁用代理信任
	//	r.SetTrustedProxies(nil)
	//}
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
		initializer.Log.Infof("初始化基础数据失败: %v", err)
	}
	initializer.Log.Info("初始化基础数据成功")
	return &App{
		app: r,
	}
}
