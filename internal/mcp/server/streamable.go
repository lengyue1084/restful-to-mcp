package server

import (
	"context"
	"flow-bridge-mcp/pkg/logger"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

type StreamableHttpTransprot struct {
	StreamableTransport transport.ServerTransport
	StreamableHandler   *transport.StreamableHTTPHandler
}

func NewMcpTransport(log *logger.Logger) (*StreamableHttpTransprot, func(), error) {
	ctx := context.Background()
	streamableTransport, handler, err := transport.NewStreamableHTTPServerTransportAndHandler()
	if err != nil {
		log.WithContext(ctx).Errorf("Failed to create streamable server: %v", err)
		return nil, nil, err
	}
	var clean = func() {
		err := streamableTransport.Shutdown(ctx, ctx)
		if err != nil {
			log.WithContext(ctx).Errorf("Failed to shutdown streamable server: %v", err)
			return
		}
	}
	return &StreamableHttpTransprot{
		StreamableTransport: streamableTransport,
		StreamableHandler:   handler,
	}, clean, nil

}
