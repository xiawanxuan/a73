package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ocean-survey-backend/models"
)

type TerrainLabelService struct {
	DB *gorm.DB
}

type CreateTerrainLabelDTO struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type UpdateTerrainLabelDTO struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func NewTerrainLabelService(db *gorm.DB) *TerrainLabelService {
	return &TerrainLabelService{DB: db}
}

func (s *TerrainLabelService) List(userID *uuid.UUID) ([]models.TerrainLabel, error) {
	var list []models.TerrainLabel
	query := s.DB.Where("is_system = ?", true)
	if userID != nil && *userID != uuid.Nil {
		query = query.Or("user_id = ?", *userID)
	}
	if err := query.Order("is_system DESC, created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TerrainLabelService) GetByID(id uuid.UUID) (*models.TerrainLabel, error) {
	var label models.TerrainLabel
	if err := s.DB.First(&label, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("terrain label not found")
		}
		return nil, err
	}
	return &label, nil
}

func (s *TerrainLabelService) Create(userID uuid.UUID, dto *CreateTerrainLabelDTO) (*models.TerrainLabel, error) {
	var existing int64
	s.DB.Model(&models.TerrainLabel{}).
		Where("name = ? AND user_id = ?", dto.Name, userID).
		Where("deleted_at IS NULL").
		Count(&existing)
	if existing > 0 {
		return nil, errors.New("label name already exists for this user")
	}

	uid := userID
	label := &models.TerrainLabel{
		UserID:      &uid,
		Name:        dto.Name,
		Color:       dto.Color,
		Description: dto.Description,
		Icon:        dto.Icon,
		IsSystem:    false,
	}

	if err := s.DB.Create(label).Error; err != nil {
		return nil, err
	}

	return s.GetByID(label.ID)
}

func (s *TerrainLabelService) Update(id uuid.UUID, userID uuid.UUID, dto *UpdateTerrainLabelDTO) (*models.TerrainLabel, error) {
	label, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if label.IsSystem {
		return nil, errors.New("cannot modify system label")
	}
	if label.UserID == nil || *label.UserID != userID {
		return nil, errors.New("permission denied: label does not belong to user")
	}

	label.Name = dto.Name
	label.Color = dto.Color
	label.Description = dto.Description
	label.Icon = dto.Icon

	if err := s.DB.Save(label).Error; err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *TerrainLabelService) Delete(id uuid.UUID, userID uuid.UUID) error {
	label, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if label.IsSystem {
		return errors.New("cannot delete system label")
	}
	if label.UserID == nil || *label.UserID != userID {
		return errors.New("permission denied: label does not belong to user")
	}

	var annotationCount int64
	s.DB.Model(&models.Annotation{}).
		Where("label_id = ?", id).
		Count(&annotationCount)
	if annotationCount > 0 {
		return errors.New("cannot delete label: it is used by annotations")
	}

	result := s.DB.Delete(label)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("terrain label not found")
	}
	return nil
}

func (s *TerrainLabelService) IsSystemLabel(id uuid.UUID) (bool, error) {
	label, err := s.GetByID(id)
	if err != nil {
		return false, err
	}
	return label.IsSystem, nil
}
