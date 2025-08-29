package api

import (
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/pkg/gormtype"
	"flow-bridge-mcp/pkg/const"
)

type OpenapiUploadRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=200"`
	FileContent string `json:"file_content" binding:"required"`
	Description string `json:"description"`
	Suffix      string `json:"suffix" binding:"required"`
}

type OpenapiUploadResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name" yaml:"name"`
	UUID        string              `json:"uuid" yaml:"uuid"`
	Description string              `json:"description" yaml:"description"`
	Urls        []string            `json:"urls,omitempty" yaml:"urls,omitempty"`
	CreatedAt   *gormtype.LocalTime `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   *gormtype.LocalTime `json:"updatedAt" yaml:"updatedAt"`
	Tools       []*ToolInfo         `json:"tools,omitempty" yaml:"tools,omitempty"`
	Version     string              `json:"version"`
	AllTools    []string            `json:"allTools"`
	Status      _const.ServerStatus `json:"status"`
}

type ToolInfo struct {
	ID             uint                       `json:"ID"`
	McpServerId    int64                      `json:"McpServerId"`
	UUID           string                     `json:"uuid"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	McpServerType  _const.McpServerTypeStatus `json:"McpServerType"`
	Method         string                     `json:"method"`
	Endpoint       string                     `json:"endpoint"`
	Headers        string                     `json:"headers"`
	Args           string                     `json:"args"`
	RequestBody    string                     `json:"requestBody"`
	ResponseBody   string                     `json:"responseBody"`
	ToolSchema     *config.ToolSchema         `json:"toolSchema"`
	Annotations    string                     `json:"annotations"`
	Security       string                     `json:"security"`
	IsAuth         _const.AuthStatus          `json:"isAuth"`
	AuthMode       string                     `json:"authMode"`
	IsPlatformAuth _const.AuthStatus          `json:"isPlatformAuth"`
	IsShow         _const.Status              `json:"isShow"`
	CreatedAt      *gormtype.LocalTime        `json:"createdAt"`
	UpdatedAt      *gormtype.LocalTime        `json:"updatedAt"`
}

type OpenapiUpdateForAuthRequest struct {
	UUID          string            `json:"uuid" binding:"required"`
	IsAuth        _const.AuthStatus `json:"isAuth" binding:"required"` //是否授权状态，这个状态是针对平台授权
	ServiceToken  string            `json:"serviceToken" binding:"omitempty,min=1,max=200"`
	PlatformToken string            `json:"platformToken" binding:"omitempty,min=1,max=200"`
	Tools         []*Tools          `json:"tools" binding:"required"`
}
type OpenapiUpdateForAuthResponse struct {
}
type Tools struct {
	ID     uint              `json:"id" binding:"required"`
	IsAuth _const.AuthStatus `json:"isAuth" binding:"required"`
}
