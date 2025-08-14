package cache

import (
	"context"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(config *conf.Conf, log *logger.Logger) (*redis.Client, func(), error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Conf.GetString("data.redis.addr"), config.Conf.GetString("data.redis.port")),
		Username:     config.Conf.GetString("data.redis.username"),
		Password:     config.Conf.GetString("data.redis.password"),
		DB:           config.Conf.GetInt("data.redis.default_db"), // use default DB
		DialTimeout:  config.Conf.GetDuration("data.redis.dial_timeout"),
		ReadTimeout:  config.Conf.GetDuration("data.redis.read_timeout"),
		WriteTimeout: config.Conf.GetDuration("data.redis.write_timeout"),
		PoolSize:     config.Conf.GetInt("data.redis.pool_size"),
		MinIdleConns: config.Conf.GetInt("data.redis.min_idle_conns"),
		MaxRetries:   config.Conf.GetInt("data.redis.max_retries"),
	})
	_, err := rdb.Ping(ctx).Result()
	cleanup := func() {
		if err := rdb.Close(); err != nil {
			log.Errorf("close redis error:%+v", err)
		}
		log.Infof("close redis success")
	}
	if err != nil {
		log.Error("open redis error:%+v", err)
		return nil, cleanup, err
	}

	log.Infof("open redis success")
	return rdb, cleanup, nil
}
