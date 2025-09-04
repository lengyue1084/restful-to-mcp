package proxy

import (
	"context"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/google/wire"
	"reflect"
	"time"
)

var ProviderSet = wire.NewSet(
	NewHttpProxy,
)

type HttpProxy struct {
	log   *logger.Logger
	cache *cache.MemoryCache
}

func NewHttpProxy(log *logger.Logger, cache *cache.MemoryCache) *HttpProxy {
	return &HttpProxy{
		log:   log,
		cache: cache,
	}
}

type ArgStruct struct {
	Name  string
	Value interface{}
	Type  string
}

func (h *HttpProxy) ctToArgs(ctMap map[string]interface{}) []*ArgStruct {
	var args []*ArgStruct
	for key, value := range ctMap {
		ref := reflect.TypeOf(value)
		args = append(args, &ArgStruct{
			Name:  key,
			Value: value,
			Type:  ref.String(),
		})
	}
	return args
}

func (h *HttpProxy) HandleHttpProxy(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {

	if serviceToken, ok := ctx.Value(_const.ServiceToken).(string); ok {
		h.log.Infof("Service Token: %s", serviceToken)
	}
	name := req.Name
	args := h.ctToArgs(req.Arguments)

	fmt.Printf("name: %v", name)
	fmt.Printf("args: %v", args)
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("无效的时区: %v", err)
	}

	return &protocol.CallToolResult{
		IsError: false,
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: time.Now().In(loc).String(),
			},
		},
	}, nil
}
