package middleware

import (
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var ProviderSet = wire.NewSet(NewMiddleware)

type middleware interface {
	Cors() gin.HandlerFunc
	Logger() gin.HandlerFunc
	Recovery() gin.HandlerFunc
	ZapLogger(logger *conf.Logger) gin.HandlerFunc
	TraceId() gin.HandlerFunc
}

type Middleware struct{}

var _ middleware = (*Middleware)(nil)

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func (m *Middleware) Logger() gin.HandlerFunc {
	return gin.Logger()
}

func (m *Middleware) Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

func (m *Middleware) ZapLogger(logger *conf.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		traceId := c.GetString("traceId")

		// 创建带有请求上下文的logger
		requestLogger := logger.With(
			zap.String("trace_id", traceId),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
		)
		// 替换Gin的默认logger
		c.Set("logger", requestLogger)
		// 处理请求
		c.Next()
		// 记录响应
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		fields := []zap.Field{
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		}
		if status >= 500 {
			requestLogger.Error("server error", fields...)
		} else if status >= 400 {
			requestLogger.Warn("client error", fields...)
		} else {
			requestLogger.Info("request completed", fields...)
		}
	}
}

func (m *Middleware) TraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceId")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		c.Set("traceId", traceId)
		c.Writer.Header().Set("traceId", traceId)
		c.Next()
	}
}
