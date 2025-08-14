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
