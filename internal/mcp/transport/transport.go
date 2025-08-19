package transport

import (
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/mcp/context"
)

//
//// 基础工具接口，提供工具信息
//type BaseTool interface {
//	Info(ctx context.Context) (*schema.ToolInfo, error)
//}
//
//// 可调用的工具接口，支持同步调用
//type InvokableTool interface {
//	BaseTool
//	InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
//}
//
//// 支持流式输出的工具接口
//type StreamableTool interface {

// TransportType represents the type of transport
type TransportType string

const (
	TypeSSE TransportType = "sse"
	//TypeStdio TransportType = "stdio"
	TypeStreamable TransportType = "streamable-http"
)

// Transport defines the interface for MCP transport implementations
type Transport interface {
	FetchTools(ctx context.Context) ([]config.ToolSchema, error)
	CallTool(ctx context.Context, params config.CallToolParams, req *context.RequestWrapper) (*config.CallToolResult, error)

	Start(ctx context.Context, tmplCtx *context.Context) error

	// Stop stops the transport
	Stop(ctx context.Context) error

	// IsRunning returns true if the transport is running
	IsRunning() bool

	// FetchPrompts fetches the list of available prompts
	FetchPrompts(ctx context.Context) ([]config.PromptSchema, error)
	// FetchPrompt fetches a specific prompt by name
	FetchPrompt(ctx context.Context, name string) (*config.PromptSchema, error)
}

// NewTransport creates transport based on the configuration
//func NewTransport(cfg config.MCPServerConfig) (Transport, error) {
//	switch TransportType(cfg.Type) {
//	case TypeSSE:
//		return &SSETransport{cfg: cfg}, nil
//	case TypeStdio:
//		return &StdioTransport{cfg: cfg}, nil
//	case TypeStreamable:
//		return &StreamableTransport{cfg: cfg}, nil
//	default:
//		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
//	}
//}
