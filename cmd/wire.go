//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data"
	"flow-bridge-mcp/internal/data/database"
	"flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/internal/mcp/transformer/openapi"
	"flow-bridge-mcp/internal/pkg/cache"
	"flow-bridge-mcp/internal/service"
	"flow-bridge-mcp/middleware"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/router"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
	))
}
