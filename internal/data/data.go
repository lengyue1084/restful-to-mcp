package data

import (
	"flow-bridge-mcp/internal/pkg/cache"
	"flow-bridge-mcp/internal/pkg/database"
	"flow-bridge-mcp/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	//NewMysqlClient,
	cache.NewRedisClient,
	database.NewPgClient,
	NewOpenapiRepo,
)

type Data struct {
	db    *gorm.DB
	redis *redis.Client
	log   logger.Logger
}

func NewData(db *gorm.DB, redis *redis.Client, log *logger.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{db: db, redis: redis}, cleanup, nil
}
