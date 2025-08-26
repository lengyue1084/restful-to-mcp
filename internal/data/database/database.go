package database

import (
	"flow-bridge-mcp/internal/pkg/cache"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	//NewMysqlClient,
	cache.NewRedisClient,
	NewPgClient,
	NewData,
	NewTransaction,
)
