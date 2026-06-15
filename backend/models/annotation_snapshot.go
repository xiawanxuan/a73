package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationCreate   OperationType = "create"
	OperationUpdate   OperationType = "update"
	OperationDelete   OperationType = "delete"
	OperationRollback OperationType = "rollback"
)

type SnapshotJSON map[string]interface{}

func (s SnapshotJSON) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SnapshotJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal SnapshotJSON value: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

type AnnotationSnapshot struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AnnotationID uuid.UUID     `gorm:"column:annotation_id;type:uuid;not null;index" json:"annotationId"`
	PointCloudID uuid.UUID     `gorm:"column:point_cloud_id;type:uuid;not null;index" json:"pointCloudId"`
	Version      int           `gorm:"column:version;not null;uniqueIndex:idx_annotation_version" json:"version"`
	Snapshot     SnapshotJSON  `gorm:"column:snapshot_json;type:jsonb;not null" json:"snapshot"`
	OperatorID   uuid.UUID     `gorm:"column:operator_id;type:uuid;not null;index" json:"operatorId"`
	Operation    OperationType `gorm:"column:operation;type:varchar(32);not null" json:"operation"`
	CreatedAt    time.Time     `gorm:"column:created_at;not null;default:current_timestamp" json:"createdAt"`
}

func (AnnotationSnapshot) TableName() string { return "annotation_snapshots" }
