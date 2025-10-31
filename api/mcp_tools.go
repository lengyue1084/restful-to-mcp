package api

import (
	"encoding/json"
	"flow-bridge-mcp/internal/mcp/config"
	_const "flow-bridge-mcp/pkg/const"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type GetMcpServerToolsRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpServerToolsResponse struct {
	Tools []*ToolProtocolInfo `json:"tools"`
}
type ToolProtocolInfo struct {
	ID             uint              `json:"id"`
	UUID           string            `json:"uuid"`
	IsAuth         _const.AuthStatus `json:"isAuth"`
	AuthMode       string            `json:"authMode"`
	IsPlatformAuth _const.AuthStatus `json:"isPlatformAuth"`
	IsShow         _const.Status     `json:"isShow"`
	// Name is the unique identifier of the tool
	Name           string                    `json:"name"`
	Description    string                    `json:"description,omitempty"`
	InputSchema    protocol.InputSchema      `json:"inputSchema"`
	OutputSchema   protocol.OutputSchema     `json:"outputSchema"`
	Annotations    *protocol.ToolAnnotations `json:"annotations,omitempty"`
	RawInputSchema json.RawMessage           `json:"-"`
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
	McpServerUUID  string             `json:"mcpServerUUID" binding:"required"`
	Name           string             `json:"name" binding:"required"`
	Description    string             `json:"description" binding:"required"`
	Method         string             `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	Path           string             `json:"path" binding:"required"`
	IsPlatformAuth _const.AuthStatus  `json:"isAuth" binding:"required"`
	IsAuth         _const.AuthStatus  `json:"isPlatformAuth" binding:"required"`
	Args           []config.ArgConfig `json:"args"`
}

type CreateMcpServerToolResponse struct {
	UUID string `json:"uuid"`
}

type UpdateMcpServerToolRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	Method      string `json:"method" binding:"required,oneof=GET POST PUT DELETE"` //binding:"oneof= GET POST PUT DELETE"
	Path        string `json:"path" binding:"required"`
	//IsShow         _const.Status      `json:"isShow" binding:"required"`
	IsPlatformAuth _const.AuthStatus  `json:"isPlatformAuth" binding:"required"`
	IsAuth         _const.AuthStatus  `json:"isAuth" binding:"required"`
	Args           []config.ArgConfig `json:"args"`
}

func ValidMethods(method string) error {
	if method == "" {
		return fmt.Errorf("method is required")
	}
	var methods = map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}
	if _, ok := methods[method]; !ok {
		return fmt.Errorf("method is not valid")
	}
	return nil
}

func ValidAuthMode(authMode config.AuthMode) error {
	if authMode == "" {
		return nil
	}
	var authModeMap = map[config.AuthMode]bool{
		config.AuthModeApiKey: true,
		config.AuthModeHttp:   true,
	}
	if _, ok := authModeMap[authMode]; !ok {
		return fmt.Errorf("authMode is not valid")
	}
	return nil
}

type UpdateMcpServerToolResponse struct {
}

type GetToolsInfoByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}
type GetToolsInfoByUUIDResponse struct {
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
	//Endpoint      string                     `json:"endpoint"`
	BaseUrl string             `json:"baseUrl"`
	Path    string             `json:"path"`
	Headers map[string]string  `json:"headers"`
	Args    []config.ArgConfig `json:"args"`
	//RequestBody   string                     `json:"requestBody"`
	//ResponseBody  string                     `json:"responseBody"`
	//ToolSchema     string                     `json:"toolSchema"`
	//Annotations    string                     `json:"annotations"`
	Security       config.Security     `json:"security"`
	IsAuth         _const.AuthStatus   `json:"isAuth"`
	AuthMode       string              `json:"authMode"`
	IsPlatformAuth _const.AuthStatus   `json:"isPlatformAuth"`
	IsShow         _const.Status       `json:"isShow"`
	SerialNumber   string              `json:"serialNumber"`
	IsRepeat       _const.CommonStatus `json:"isRepeat"`
}

type TestMcpServerToolRequest struct {
}
type TestMcpServerToolResponse struct {
}
