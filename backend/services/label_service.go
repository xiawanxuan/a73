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

func (s *TerrainLabelService) List() ([]models.TerrainLabel, error) {
	var list []models.TerrainLabel
	if err := s.DB.Find(&list).Error; err != nil {
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

func (s *TerrainLabelService) Create(dto *CreateTerrainLabelDTO) (*models.TerrainLabel, error) {
	label := &models.TerrainLabel{
		Name:        dto.Name,
		Color:       dto.Color,
		Description: dto.Description,
		Icon:        dto.Icon,
	}

	if err := s.DB.Create(label).Error; err != nil {
		return nil, err
	}

	return s.GetByID(label.ID)
}

func (s *TerrainLabelService) Update(id uuid.UUID, dto *UpdateTerrainLabelDTO) (*models.TerrainLabel, error) {
	label, err := s.GetByID(id)
	if err != nil {
		return nil, err
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

func (s *TerrainLabelService) Delete(id uuid.UUID) error {
	result := s.DB.Delete(&models.TerrainLabel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("terrain label not found")
	}
	return nil
}
