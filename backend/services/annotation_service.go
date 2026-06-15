package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ocean-survey-backend/models"
)

type AnnotationService struct {
	DB *gorm.DB
}

type AnnotationWithLabel struct {
	models.Annotation
	Label *models.TerrainLabel `json:"label,omitempty"`
}

type CreateAnnotationDTO struct {
	Name    string             `json:"name"`
	Polygon models.PolygonJSON `json:"polygon"`
	LabelID string             `json:"labelId"`
}

type UpdateAnnotationDTO struct {
	Name    string             `json:"name"`
	Polygon models.PolygonJSON `json:"polygon"`
	LabelID string             `json:"labelId"`
}

func NewAnnotationService(db *gorm.DB) *AnnotationService {
	return &AnnotationService{DB: db}
}

func (s *AnnotationService) ListByPointCloud(pointCloudID uuid.UUID) ([]AnnotationWithLabel, error) {
	var annotations []models.Annotation
	if err := s.DB.Where("point_cloud_id = ? AND deleted_at IS NULL", pointCloudID).Find(&annotations).Error; err != nil {
		return nil, err
	}

	result := make([]AnnotationWithLabel, 0, len(annotations))
	for _, ann := range annotations {
		aw := AnnotationWithLabel{Annotation: ann}
		var label models.TerrainLabel
		if err := s.DB.First(&label, "id = ?", ann.LabelID).Error; err == nil {
			aw.Label = &label
		}
		result = append(result, aw)
	}

	return result, nil
}

func (s *AnnotationService) GetByID(id uuid.UUID) (*AnnotationWithLabel, error) {
	var ann models.Annotation
	if err := s.DB.First(&ann, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("annotation not found")
		}
		return nil, err
	}

	aw := &AnnotationWithLabel{Annotation: ann}
	var label models.TerrainLabel
	if err := s.DB.First(&label, "id = ?", ann.LabelID).Error; err == nil {
		aw.Label = &label
	}

	return aw, nil
}

func (s *AnnotationService) validateAnnotation(pointCloudID uuid.UUID, dto *CreateAnnotationDTO) error {
	if dto.Polygon == nil || len(dto.Polygon) < 3 {
		return errors.New("polygon must have at least 3 points")
	}
	labelID, err := uuid.Parse(dto.LabelID)
	if err != nil {
		return errors.New("invalid label id")
	}
	var labelCount int64
	if err := s.DB.Model(&models.TerrainLabel{}).Where("id = ?", labelID).Count(&labelCount).Error; err != nil {
		return err
	}
	if labelCount == 0 {
		return errors.New("label not found")
	}
	if pointCloudID == uuid.Nil {
		return errors.New("point cloud id is required")
	}
	var pcCount int64
	if err := s.DB.Model(&models.PointCloud{}).Where("id = ?", pointCloudID).Count(&pcCount).Error; err != nil {
		return err
	}
	if pcCount == 0 {
		return errors.New("point cloud not found")
	}
	return nil
}

func (s *AnnotationService) computeBoundsCenter(polygon models.PolygonJSON) (cx, cy, cz float64) {
	if len(polygon) == 0 {
		return 0, 0, 0
	}
	var sx, sy, sz float64
	for _, p := range polygon {
		sx += p.X
		sy += p.Y
		sz += p.Z
	}
	n := float64(len(polygon))
	return sx / n, sy / n, sz / n
}

func (s *AnnotationService) buildSnapshotJSON(ann *models.Annotation) models.SnapshotJSON {
	return models.SnapshotJSON{
		"id":              ann.ID.String(),
		"point_cloud_id":  ann.PointCloudID.String(),
		"label_id":        ann.LabelID.String(),
		"name":            ann.Name,
		"polygon_json":    ann.Polygon,
		"bounds_center_x": ann.BoundsCenterX,
		"bounds_center_y": ann.BoundsCenterY,
		"bounds_center_z": ann.BoundsCenterZ,
		"creator_id":      ann.CreatorID.String(),
	}
}

func (s *AnnotationService) getNextVersion(annotationID uuid.UUID) (int, error) {
	var maxVersion int
	err := s.DB.Model(&models.AnnotationSnapshot{}).
		Where("annotation_id = ?", annotationID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (s *AnnotationService) createSnapshot(tx *gorm.DB, ann *models.Annotation, operatorID uuid.UUID, operation models.OperationType, version int) error {
	snapshot := &models.AnnotationSnapshot{
		ID:           uuid.New(),
		AnnotationID: ann.ID,
		PointCloudID: ann.PointCloudID,
		Version:      version,
		Snapshot:     s.buildSnapshotJSON(ann),
		OperatorID:   operatorID,
		Operation:    operation,
		CreatedAt:    time.Now(),
	}
	return tx.Create(snapshot).Error
}

func (s *AnnotationService) ValidateAndCreate(pointCloudID uuid.UUID, dto *CreateAnnotationDTO, userID uuid.UUID) (*AnnotationWithLabel, error) {
	if err := s.validateAnnotation(pointCloudID, dto); err != nil {
		return nil, err
	}

	labelID, _ := uuid.Parse(dto.LabelID)
	cx, cy, cz := s.computeBoundsCenter(dto.Polygon)

	ann := &models.Annotation{
		ID:            uuid.New(),
		PointCloudID:  pointCloudID,
		LabelID:       labelID,
		Name:          dto.Name,
		Polygon:       dto.Polygon,
		BoundsCenterX: cx,
		BoundsCenterY: cy,
		BoundsCenterZ: cz,
		CreatorID:     userID,
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ann).Error; err != nil {
			return err
		}
		if err := s.createSnapshot(tx, ann, userID, models.OperationCreate, 1); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(ann.ID)
}

func (s *AnnotationService) ValidateAndUpdate(id uuid.UUID, dto *UpdateAnnotationDTO, userID uuid.UUID) (*AnnotationWithLabel, error) {
	ann, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	labelID, labelErr := uuid.Parse(dto.LabelID)
	if labelErr != nil {
		return nil, errors.New("invalid label id")
	}

	createDTO := &CreateAnnotationDTO{
		Name:    dto.Name,
		Polygon: dto.Polygon,
		LabelID: dto.LabelID,
	}
	if err := s.validateAnnotation(ann.PointCloudID, createDTO); err != nil {
		return nil, err
	}

	cx, cy, cz := s.computeBoundsCenter(dto.Polygon)

	ann.Name = dto.Name
	ann.LabelID = labelID
	ann.Polygon = dto.Polygon
	ann.BoundsCenterX = cx
	ann.BoundsCenterY = cy
	ann.BoundsCenterZ = cz

	nextVersion, err := s.getNextVersion(id)
	if err != nil {
		return nil, err
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&ann.Annotation).Error; err != nil {
			return err
		}
		if err := s.createSnapshot(tx, &ann.Annotation, userID, models.OperationUpdate, nextVersion); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *AnnotationService) Delete(id uuid.UUID, userID uuid.UUID) error {
	var ann models.Annotation
	if err := s.DB.First(&ann, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("annotation not found")
		}
		return err
	}

	nextVersion, err := s.getNextVersion(id)
	if err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&ann).Error; err != nil {
			return err
		}
		if err := s.createSnapshot(tx, &ann, userID, models.OperationDelete, nextVersion); err != nil {
			return err
		}
		return nil
	})
}

func (s *AnnotationService) Rollback(annotationID uuid.UUID, version int, userID uuid.UUID) (*AnnotationWithLabel, error) {
	var snapshot models.AnnotationSnapshot
	if err := s.DB.Where("annotation_id = ? AND version = ?", annotationID, version).
		First(&snapshot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("snapshot not found")
		}
		return nil, err
	}

	var ann models.Annotation
	if err := s.DB.Unscoped().First(&ann, "id = ?", annotationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("annotation not found")
		}
		return nil, err
	}

	sj := snapshot.Snapshot
	if v, ok := sj["label_id"].(string); ok {
		parsed, err := uuid.Parse(v)
		if err == nil {
			ann.LabelID = parsed
		}
	}
	if v, ok := sj["name"].(string); ok {
		ann.Name = v
	}
	if v, ok := sj["polygon_json"]; ok {
		pjBytes, _ := json.Marshal(v)
		var pj models.PolygonJSON
		json.Unmarshal(pjBytes, &pj)
		ann.Polygon = pj
	}
	if v, ok := sj["bounds_center_x"].(float64); ok {
		ann.BoundsCenterX = v
	}
	if v, ok := sj["bounds_center_y"].(float64); ok {
		ann.BoundsCenterY = v
	}
	if v, ok := sj["bounds_center_z"].(float64); ok {
		ann.BoundsCenterZ = v
	}

	ann.DeletedAt = gorm.DeletedAt{}

	nextVersion, err := s.getNextVersion(annotationID)
	if err != nil {
		return nil, err
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Save(&ann).Error; err != nil {
			return err
		}
		if err := s.createSnapshot(tx, &ann, userID, models.OperationRollback, nextVersion); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(annotationID)
}

func (s *AnnotationService) ListSnapshots(annotationID uuid.UUID) ([]models.AnnotationSnapshot, error) {
	var snapshots []models.AnnotationSnapshot
	if err := s.DB.Where("annotation_id = ?", annotationID).
		Order("version DESC").
		Limit(30).
		Find(&snapshots).Error; err != nil {
		return nil, err
	}
	return snapshots, nil
}
