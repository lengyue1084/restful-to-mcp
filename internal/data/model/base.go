package model

import (
	"flow-bridge-mcp/internal/pkg/gormtype"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint                `gorm:"primarykey" json:"id"`
	CreatedAt *gormtype.LocalTime `json:"created_at"`
	UpdatedAt *gormtype.LocalTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
}

type (
	AuthStatus          int8
	HaveToolsStatus     int8
	McpServerTypeStatus int8
)

const (
	IsAuthNo             AuthStatus          = 1
	IsAuthYes            AuthStatus          = 2
	HaveToolsNo          HaveToolsStatus     = 1
	HaveToolsYes         HaveToolsStatus     = 2
	McpServerTypeOpenapi McpServerTypeStatus = 1
	McpServerTypeGrpc    McpServerTypeStatus = 2
)

func (a AuthStatus) String() string {
	switch a {
	case IsAuthNo:
		return "否"
	case IsAuthYes:
		return "是"
	default:
		return "未知"
	}
}
func (h HaveToolsStatus) String() string {
	switch h {
	case HaveToolsNo:
		return "否"
	case HaveToolsYes:
		return "是"
	default:
		return "未知"
	}

}

func (m McpServerTypeStatus) String() string {
	switch m {
	case McpServerTypeOpenapi:
		return "OpenAPI"
	case McpServerTypeGrpc:
		return "gRPC"
	default:
		return "未知"
	}
}
