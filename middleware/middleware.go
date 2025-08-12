package middleware

import (
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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

func NewMiddleware() *Middleware {
	return &Middleware{}
}
