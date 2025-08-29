package api

import "flow-bridge-mcp/internal/mcp/config"

type GetMcpServerToolsRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpServerToolsResponse struct {
	Tools []*config.ToolSchema `json:"tools"`
}

type ToolItemInfo struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	InputSchema config.ToolInputSchema  `json:"inputSchema"`
	Annotations *config.ToolAnnotations `json:"annotations,omitempty"`
}
