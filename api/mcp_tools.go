package api

import (
	"flow-bridge-mcp/internal/mcp/config"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type GetMcpServerToolsRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpServerToolsResponse struct {
	Tools []*protocol.Tool `json:"tools"`
}

type ToolItemInfo struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	InputSchema config.ToolInputSchema  `json:"inputSchema"`
	Annotations *config.ToolAnnotations `json:"annotations,omitempty"`
}
