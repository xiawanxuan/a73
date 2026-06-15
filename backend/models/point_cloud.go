package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PointCloud struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	FilePath    string    `gorm:"column:file_path;type:varchar(512);not null" json:"filePath"`
	FileSize    int64     `gorm:"column:file_size;not null;default:0" json:"fileSize"`
	PointCount  int64     `gorm:"column:point_count;not null;default:0" json:"pointCount"`
	BoundsMinX  float64   `gorm:"column:bounds_min_x;not null;default:0" json:"boundsMinX"`
	BoundsMinY  float64   `gorm:"column:bounds_min_y;not null;default:0" json:"boundsMinY"`
	BoundsMinZ  float64   `gorm:"column:bounds_min_z;not null;default:0" json:"boundsMinZ"`
	BoundsMaxX  float64   `gorm:"column:bounds_max_x;not null;default:0" json:"boundsMaxX"`
	BoundsMaxY  float64   `gorm:"column:bounds_max_y;not null;default:0" json:"boundsMaxY"`
	BoundsMaxZ  float64   `gorm:"column:bounds_max_z;not null;default:0" json:"boundsMaxZ"`
	UploadedBy  uuid.UUID `gorm:"column:uploaded_by;type:uuid;not null;index" json:"uploadedBy"`
	Status      string    `gorm:"column:status;type:varchar(32);not null;default:'processing'" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:current_timestamp" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (PointCloud) TableName() string { return "point_clouds" }
