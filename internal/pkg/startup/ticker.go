package startup

import (
	"context"
	"flow-bridge-mcp/internal/biz"
	mcpServer "flow-bridge-mcp/internal/mcp/server"
	_const "flow-bridge-mcp/pkg/const"
	"time"
)

type Ticker struct {
	mcpServerTickerStop chan bool
	openapiUserCase     *biz.OpenapiUseCase
	mcpServerManager    *mcpServer.McpServerManager
}

func NewTicker(openapiUserCase *biz.OpenapiUseCase,
	mcpServerManager *mcpServer.McpServerManager) *Ticker {
	return &Ticker{
		mcpServerTickerStop: make(chan bool),
		openapiUserCase:     openapiUserCase,
		mcpServerManager:    mcpServerManager,
	}
}

func (i *Ticker) StartMcpServerCacheRefreshTicker() {
	var timer = time.NewTicker(time.Second * _const.McpServerRefreshTicketTime)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*_const.CommonContextTimeOut)
			i.openapiUserCase.UpdateToolsForCache(ctx)
			i.mcpServerManager.RegisterToolFromCache()
			cancel()
		case <-i.mcpServerTickerStop:
			return
		}
	}
}
func (i *Ticker) StopMcpServerCacheRefreshTicker() {
	i.mcpServerTickerStop <- true
}
