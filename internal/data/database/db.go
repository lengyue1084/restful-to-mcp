package database

import (
	"context"
	"flow-bridge-mcp/pkg/logger"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {

	err = d.Db.Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context
		ctx = context.WithValue(ctx, "tx", tx)
		// 执行业务逻辑
		return fn(ctx)
		// 如果 fn 返回 error，GORM 自动回滚
		// 如果 fn 返回 nil，GORM 自动提交
	})

	return err
}
