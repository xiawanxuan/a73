package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TerrainLabel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      *uuid.UUID     `gorm:"column:user_id;type:uuid;index" json:"userId,omitempty"`
	Name        string         `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Color       string         `gorm:"column:color;type:varchar(16);not null;default:'#1E88E5'" json:"color"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Icon        string         `gorm:"column:icon;type:varchar(64)" json:"icon"`
	IsSystem    bool           `gorm:"column:is_system;not null;default:false" json:"isSystem"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;default:current_timestamp" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (TerrainLabel) TableName() string { return "terrain_labels" }
