package api

import (
	"restful-to-mcp/internal/mcp/config"
	_const "restful-to-mcp/pkg/const"
)

type GetMcpServerInfoByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type ListMcpServersRequest struct {
}

type GetMcpServerInfoByUUIDResponse struct {
	ID        uint                `json:"id"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
	Status    _const.ServerStatus `json:"status"`
	Source    _const.SourceType   `json:"source"`
	CommonMcpServerByForm
}

type McpServerListItem struct {
	ID          uint                  `json:"id"`
	UUID        string                `json:"uuid"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Version     string                `json:"version"`
	Status      _const.ServerStatus   `json:"status"`
	Source      _const.SourceType     `json:"source"`
	IsAuth      _const.AuthTypeStatus `json:"isAuth"`
	ToolCount   int                   `json:"toolCount"`
	CreatedAt   string                `json:"createdAt"`
	UpdatedAt   string                `json:"updatedAt"`
}

type ListMcpServersResponse struct {
	Items []*McpServerListItem `json:"items"`
}

type UpdateMcpServerByUUIDRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type UpdateMcpServerByUUIDResponse struct {
}

type GetMcpConnectTokenByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpConnectTokenByUUIDResponse struct {
	ConnectToken string `json:"connectToken" binding:"required"`
}

type DeleteMcpServerByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}
type DeleteMcpServerByUUIDResponse struct {
}

type CommonMcpServerByForm struct {
	UUID          string                `json:"uuid" binding:"required"`
	Name          string                `json:"name" binding:"required"`
	Description   string                `json:"description" binding:"required"`
	Urls          []string              `json:"urls" binding:"required"`
	Version       string                `json:"version" binding:"required"`
	IsAuth        _const.AuthTypeStatus `json:"isAuth" binding:"required"` //0未知，1 不开启，2开启service授权，3开启平台授权，4开启所有的授权'
	PlatformToken string                `json:"platformToken"`
	ServiceToken  string                `json:"serviceToken"`
	Headers       map[string]string     `json:"headers"`
	Security      config.Security       `json:"security"`
}

type CreateMcpServerByFormRequest struct {
	CommonMcpServerByForm
}

type CreateMcpServerByFormResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"createdAt"`
}

type UpdateMcpServerByFormRequest struct {
	CommonMcpServerByForm
}

type UpdateMcpServerByFormResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"createdAt"`
}
