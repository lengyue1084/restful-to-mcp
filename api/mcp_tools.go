package api

import (
	"flow-bridge-mcp/internal/mcp/config"
	_const "flow-bridge-mcp/pkg/const"
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

type CommonToolItemInfo struct {
	ID            uint                       `json:"id"`
	UUID          string                     `json:"uuid"`
	CreatedAt     string                     `json:"createdAt"`
	UpdatedAt     string                     `json:"updatedAt"`
	McpServerId   int64                      `json:"mcpServerId"`
	McpServerUUID string                     `json:"mcpServerUUID"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	McpServerType _const.McpServerTypeStatus `json:"mcpServerType"`
	Method        string                     `json:"method"`
	Endpoint      string                     `json:"endpoint"`
	Headers       string                     `json:"headers"`
	//Args          string                     `json:"args"`
	//RequestBody   string                     `json:"requestBody"`
	//ResponseBody  string                     `json:"responseBody"`
	//ToolSchema     string                     `json:"toolSchema"`
	//Annotations    string                     `json:"annotations"`
	Security       string              `json:"security"`
	IsAuth         _const.AuthStatus   `json:"isAuth"`
	AuthMode       string              `json:"authMode"`
	IsPlatformAuth _const.AuthStatus   `json:"isPlatformAuth"`
	IsShow         _const.Status       `json:"isShow"`
	SerialNumber   string              `json:"serialNumber"`
	IsRepeat       _const.CommonStatus `json:"isRepeat"`
}

type GetMcpServerToolsByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpServerToolsByUUIDResponse struct {
	Tools []*CommonToolItemInfo `json:"tools"`
}

type CreateMcpServerToolRequest struct {
	McpServerUUID  string              `json:"mcpServerUUID" binding:"required"`
	Name           string              `json:"name" binding:"required"`
	Description    string              `json:"description" binding:"required"`
	Method         string              `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	Path           string              `json:"path" binding:"required"`
	IsShow         _const.Status       `json:"isShow" binding:"required"`
	IsPlatformAuth _const.AuthStatus   `json:"isAuth" binding:"required"`
	IsAuth         _const.AuthStatus   `json:"isPlatformAuth" binding:"required"`
	AuthMode       config.AuthMode     `json:"authMode" binding:"oneof=apiKey http"`
	SecurityKey    string              `json:"securityKey"`
	Position       config.AuthPosition `json:"position"` // query header
	Scheme         string              `json:"scheme"`
}

type CreateMcpServerToolResponse struct {
	UUID string `json:"uuid"`
}

type UpdateMcpServerToolRequest struct {
	UUID           string              `json:"uuid" binding:"required"`
	Name           string              `json:"name" binding:"required"`
	Description    string              `json:"description" binding:"required"`
	Method         string              `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	Path           string              `json:"path" binding:"required"`
	IsShow         _const.Status       `json:"isShow" binding:"required"`
	IsPlatformAuth _const.AuthStatus   `json:"isAuth" binding:"required"`
	IsAuth         _const.AuthStatus   `json:"isPlatformAuth" binding:"required"`
	AuthMode       config.AuthMode     `json:"authMode" binding:"oneof=apiKey http"`
	SecurityKey    string              `json:"securityKey"`
	Position       config.AuthPosition `json:"position"` // query header
	Scheme         string              `json:"scheme"`
}

type UpdateMcpServerToolResponse struct {
}
