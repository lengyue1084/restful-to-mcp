package startup

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/conf"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	"flow-bridge-mcp/pkg/logger"
	"github.com/google/wire"
	"sync"
)

// ProviderSet is initializer providers.
var ProviderSet = wire.NewSet(NewInitializer, NewTicker)

type Initializer struct {
	openapiUserCase                 *biz.OpenapiUseCase
	mcpServerManager                *mcpServer.McpServerManager
	conf                            *conf.Conf
	Log                             *logger.Logger
	initialized                     bool
	mux                             sync.Mutex
	mcpServerCacheRefreshTickerStop chan bool
	ticker                          *Ticker
}

func NewInitializer(conf *conf.Conf, ticker *Ticker, log *logger.Logger, mcpServerManager *mcpServer.McpServerManager, openapiUserCase *biz.OpenapiUseCase) *Initializer {

	return &Initializer{
		conf:             conf,
		Log:              log,
		openapiUserCase:  openapiUserCase,
		initialized:      false,
		mux:              sync.Mutex{},
		mcpServerManager: mcpServerManager,
		ticker:           ticker,
	}
}
func (i *Initializer) Initialize() (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if i.initialized {
		return nil
	}
	i.Log.Info("Starting application initialization...")
	if err = i.initializeMcpServers(); err != nil { // 加载McpServer的工具
		return err
	}
	// 其他类似服务可以从这里加载
	go func() {
		defer func() {
			if r := recover(); r != nil {
				i.Log.Error("panic: %v", r)
			}
		}()
		i.ticker.StartMcpServerCacheRefreshTicker()
	}()

	i.Log.Info("Application initialization completed")
	return
}

func (i *Initializer) initializeMcpServers() (err error) {

	i.openapiUserCase.UpdateToolsForCache(context.Background())
	i.mcpServerManager.RegisterToolFromCache()
	return nil
}

func (i *Initializer) IsInitialized() bool {
	i.mux.Lock()
	defer i.mux.Unlock()
	return i.initialized
}
