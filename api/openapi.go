package api

import (
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/pkg/gormtype"
)

type OpenapiUploadRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=200"`
	FileContent string `json:"file_content" binding:"required"`
	Description string `json:"description"`
	Suffix      string `json:"suffix" binding:"required"`
}

type OpenapiUploadResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name" yaml:"name"`
	UUID        string             `json:"uuid" yaml:"uuid"`
	Description string             `json:"description" yaml:"description"`
	Urls        []string           `json:"urls,omitempty" yaml:"urls,omitempty"`
	CreatedAt   gormtype.LocalTime `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   gormtype.LocalTime `json:"updatedAt" yaml:"updatedAt"`
	Tools       []*ToolConfig      `json:"tools,omitempty" yaml:"tools,omitempty"`
	Version     string             `json:"version"`
	AllTools    []string           `json:"allTools"`
}

type ToolConfig struct {
	ID            uint                      `json:"ID"`
	McpServerId   int64                     `json:"McpServerId"`
	UUID          string                    `json:"uuid"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	McpServerType model.McpServerTypeStatus `json:"McpServerType"`
	Method        string                    `json:"method"`
	Endpoint      string                    `json:"endpoint" gorm:"column:endpoint;type:varchar(255);default:'';comment:API端点地址"`
	Headers       string                    `json:"headers"`

	InputSchema string `json:"inputSchema"`
	Annotations string `json:"annotations"`
	Security    string `json:"security"`

	IsAuth         model.AuthStatus   `json:"isAuth"`
	AuthMode       string             `json:"authMode"`
	IsPlatformAuth model.AuthStatus   `json:"isPlatformAuth"`
	IsShow         model.Status       `json:"isShow"`
	CreatedAt      gormtype.LocalTime `json:"createdAt"`
	UpdatedAt      gormtype.LocalTime `json:"updatedAt"`
}
