package api

import (
	"flow-bridge-mcp/internal/mcp/config"
	_const "flow-bridge-mcp/pkg/const"
)

type GetMcpServerInfoByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpServerInfoByUUIDResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CommonMcpServerByForm
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
	Header        []map[string]string   `json:"header"`
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
