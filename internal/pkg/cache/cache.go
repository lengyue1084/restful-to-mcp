package cache

import (
	"context"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"sync/atomic"
)

var ProviderSet = wire.NewSet(
	NewRedisClient,
	NewMemoryCache,
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
		panic(fmt.Errorf("open redis error:%+v", err))
		return nil, cleanup, err
	}

	log.Infof("open redis success")
	return rdb, cleanup, nil
}

type MemoryCache struct {
	McpServer    atomic.Value
	OldMcpServer atomic.Value
	log          *logger.Logger
}

func NewMemoryCache(log *logger.Logger) *MemoryCache {
	return &MemoryCache{
		McpServer:    atomic.Value{},
		OldMcpServer: atomic.Value{},
		log:          log,
	}
}

type TypeCache int

const (
	NewMcpValue TypeCache = iota
	OldMcpValue
)

func (m *MemoryCache) StoreMcpServer(typeCache TypeCache, mcpServerInfo []*model.McpServer) {
	switch typeCache {
	case NewMcpValue:
		m.McpServer.Store(mcpServerInfo)
	case OldMcpValue:
		m.OldMcpServer.Store(mcpServerInfo)
	default:
		return
	}
}

func (m *MemoryCache) LoadMcpServer(typeCache TypeCache) (mcpServerList []*model.McpServer, ok bool) {

	switch typeCache {
	case NewMcpValue:
		if value := m.McpServer.Load(); value != nil {
			if serverInfo, ok := value.([]*model.McpServer); ok {
				return serverInfo, ok
			}
		}
	case OldMcpValue:
		if value := m.OldMcpServer.Load(); value != nil {
			if serverInfo, ok := value.([]*model.McpServer); ok {
				return serverInfo, ok
			}
		}
	default:
		return nil, false
	}
	return nil, false
}

func (m *MemoryCache) ClearCache(typeCache TypeCache) {
	m.StoreMcpServer(typeCache, nil)

}
