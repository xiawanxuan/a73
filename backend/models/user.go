package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string         `gorm:"column:username;type:varchar(64);not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"column:role;type:varchar(32);not null;default:'user'" json:"role"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null;default:current_timestamp" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (User) TableName() string { return "users" }
