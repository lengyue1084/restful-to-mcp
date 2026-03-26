//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"restful-to-mcp/internal/biz"
	"restful-to-mcp/internal/conf"
	"restful-to-mcp/internal/data"
	"restful-to-mcp/internal/data/database"
	"restful-to-mcp/internal/mcp/proxy"
	"restful-to-mcp/internal/mcp/server"
	"restful-to-mcp/internal/mcp/transformer/openapi"
	"restful-to-mcp/internal/pkg/cache"
	"restful-to-mcp/internal/pkg/startup"
	"restful-to-mcp/internal/service"
	"restful-to-mcp/middleware"
	"restful-to-mcp/pkg/logger"
	"restful-to-mcp/router"
)

// initApp init gin application.
func initApp(config *conf.Conf) (*gin.Engine, func(), error) {
	panic(wire.Build(
		middleware.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		router.ProviderSet,
		logger.ProviderSet,
		openapi.ProviderSet,
		database.ProviderSet,
		server.ProviderSet,
		cache.ProviderSet,
		startup.ProviderSet,
		proxy.ProviderSet,
	))
}
