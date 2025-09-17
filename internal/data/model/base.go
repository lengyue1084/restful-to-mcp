package model

import (
	"flow-bridge-mcp/internal/pkg/gormtype"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint                `gorm:"primarykey" json:"id"`
	CreatedAt *gormtype.LocalTime `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt *gormtype.LocalTime `json:"updatedAt" gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt      `json:"deletedAt" gorm:"column:deleted_at;index:idx_deleted_at"`
}
