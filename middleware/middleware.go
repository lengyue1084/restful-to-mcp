package middleware

import (
	"bytes"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

var ProviderSet = wire.NewSet(NewMiddleware)

type middleware interface {
	CorsMiddleware() gin.HandlerFunc
	LoggerMiddleware() gin.HandlerFunc
	RecoveryMiddleware() gin.HandlerFunc
	TraceIdMiddleware() gin.HandlerFunc
	TimeoutMiddleware(s time.Duration) gin.HandlerFunc
}

type Middleware struct {
	log *logger.Logger
}

var _ middleware = (*Middleware)(nil)

func NewMiddleware(log *logger.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}

func (m *Middleware) CorsMiddleware() gin.HandlerFunc {
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

func (m *Middleware) timeoutResponse(c *gin.Context) {
	m.log.WithContext(c).Error("Request Timeout")
	c.String(http.StatusRequestTimeout, "timeout")
	c.Abort()
}

// TimeoutMiddleware Custom timeout middleware
func (m *Middleware) TimeoutMiddleware(timeoutDuration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(timeoutDuration),
		timeout.WithResponse(m.timeoutResponse),
	)
}
func (m *Middleware) RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// ResponseWriter包装器，用于捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 条件性读取请求体（避免记录敏感信息）
		var requestBody string
		if shouldLogBody(path) && c.Request.Method != "GET" && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 创建基础日志字段
		baseFields := []zap.Field{
			zap.String(_const.TraceId, c.GetString(_const.TraceId)),
			zap.String(_const.SpanId, c.GetString(_const.SpanId)),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("requestBody", requestBody),
		}

		// 根据需要添加请求体
		//if requestBody != "" {
		//	baseFields = append(baseFields, zap.String("request_body", requestBody))
		//}

		//requestLogger := logger.With(baseFields)
		//c.Set("logger", requestLogger)

		// 包装ResponseWriter
		var responseBody bytes.Buffer
		c.Writer = &responseBodyWriter{c.Writer, &responseBody}

		// 处理请求
		c.Next()

		// 记录响应
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		// 构建响应字段
		responseFields := []zap.Field{
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("bodySize", c.Writer.Size()),
			zap.String("responseBody", responseBody.String()),
		}

		baseFields = append(baseFields, responseFields...)

		// 条件性记录响应体
		//if shouldLogBody(path) && status >= 400 {
		//	responseFields = append(responseFields, zap.String("response_body", responseBody.String()))
		//}
		//responseFields = append(baseFields, zap.String("response_body", responseBody.String()))

		// 根据状态码记录不同级别的日志
		//if status >= 500 {
		//	m.log.Errorf("server error:%+v", baseFields)
		//} else if status >= 400 {
		//	m.log.Warnf("client error:%+v", baseFields)
		//} else {
		//	m.log.Infof("request completed:%+v", baseFields)
		//}

		logMsg := "request completed"
		switch {
		case status >= 500:
			logMsg = "server error"
			m.log.WithZapFields(baseFields...).Error(logMsg)
		case status >= 400:
			logMsg = "client error"
			m.log.WithZapFields(baseFields...).Warn(logMsg)
		default:
			m.log.WithZapFields(baseFields...).Info(logMsg)
		}

	}
}

// 判断是否应该记录请求/响应体,可以用做过滤敏感操作
func shouldLogBody(path string) bool {
	// 可以根据路径配置哪些接口需要记录body
	//sensitivePaths := []string{"/auth/login", "/user/password"}
	//for _, p := range sensitivePaths {
	//	if strings.Contains(path, p) {
	//		return false // 敏感路径不记录body
	//	}
	//}
	return true
}

func (m *Middleware) TraceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceId")
		spanId := uuid.New().String()
		c.Set("spanId", spanId) // 为当前请求生成新的spanId
		if traceId == "" {
			traceId = spanId // 如果没有traceId，则使用spanId作为traceId
		}
		c.Set("traceId", traceId)
		c.Writer.Header().Set("traceId", traceId)
		c.Next()
	}
}
