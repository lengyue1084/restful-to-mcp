package model

import (
	"gorm.io/gorm"
	"restful-to-mcp/internal/pkg/gormtype"
)

type BaseModel struct {
	ID        uint                `gorm:"primarykey" json:"id"`
	CreatedAt *gormtype.LocalTime `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt *gormtype.LocalTime `json:"updatedAt" gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt      `json:"deletedAt" gorm:"column:deleted_at;index:idx_deleted_at"`
}
