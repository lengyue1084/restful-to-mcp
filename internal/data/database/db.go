package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"restful-to-mcp/pkg/const"
	"restful-to-mcp/pkg/logger"
)

type Data struct {
	Db    *gorm.DB
	redis *redis.Client
	log   logger.Logger
}

func NewData(db *gorm.DB, redis *redis.Client, log *logger.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{Db: db, redis: redis}, cleanup, nil
}

// GetDb biz 层开启事务可以使用该方法获取同一个tx对象，避免多表操作时，tx对象不一致
func (m *Data) GetDb(ctx context.Context) (db *gorm.DB, err error) {
	if tx, ok := ctx.Value(_const.TxKey).(*gorm.DB); ok {
		return tx, nil
	}
	return nil, fmt.Errorf("tx is not exist")
}
