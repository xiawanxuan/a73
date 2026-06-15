package models

import (
	"time"

	"github.com/google/uuid"
)

type TerrainLabel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex" json:"name"`
	Color       string    `gorm:"column:color;type:varchar(16);not null;default:'#1E88E5'" json:"color"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Icon        string    `gorm:"column:icon;type:varchar(64)" json:"icon"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:current_timestamp" json:"updatedAt"`
}

func (TerrainLabel) TableName() string { return "terrain_labels" }
