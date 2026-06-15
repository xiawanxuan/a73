package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type PolygonJSON []Point3D

func (p PolygonJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PolygonJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PolygonJSON value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}

type Annotation struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PointCloudID  uuid.UUID      `gorm:"column:point_cloud_id;type:uuid;not null;index" json:"pointCloudId"`
	LabelID       uuid.UUID      `gorm:"column:label_id;type:uuid;not null;index" json:"labelId"`
	Label         *TerrainLabel  `gorm:"foreignKey:LabelID" json:"label,omitempty"`
	Name          string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Polygon       PolygonJSON    `gorm:"column:polygon_json;type:jsonb;not null" json:"polygon"`
	BoundsCenterX float64        `gorm:"column:bounds_center_x;type:double precision;not null;default:0" json:"boundsCenterX"`
	BoundsCenterY float64        `gorm:"column:bounds_center_y;type:double precision;not null;default:0" json:"boundsCenterY"`
	BoundsCenterZ float64        `gorm:"column:bounds_center_z;type:double precision;not null;default:0" json:"boundsCenterZ"`
	CreatorID     uuid.UUID      `gorm:"column:creator_id;type:uuid;not null;index" json:"creatorId"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null;default:current_timestamp" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Annotation) TableName() string { return "annotations" }
