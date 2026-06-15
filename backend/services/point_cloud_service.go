package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ocean-survey-backend/models"
)

type PointCloudService struct {
	DB *gorm.DB
}

type CreatePointCloudDTO struct {
	Name       string    `json:"name"`
	FilePath   string    `json:"filePath"`
	FileSize   int64     `json:"fileSize"`
	PointCount int64     `json:"pointCount"`
	BoundsMinX float64   `json:"boundsMinX"`
	BoundsMinY float64   `json:"boundsMinY"`
	BoundsMinZ float64   `json:"boundsMinZ"`
	BoundsMaxX float64   `json:"boundsMaxX"`
	BoundsMaxY float64   `json:"boundsMaxY"`
	BoundsMaxZ float64   `json:"boundsMaxZ"`
	UploadedBy uuid.UUID `json:"uploadedBy"`
}

type UpdatePointCloudDTO struct {
	Name       string `json:"name"`
	FilePath   string `json:"filePath"`
	FileSize   int64  `json:"fileSize"`
	PointCount int64  `json:"pointCount"`
}

func NewPointCloudService(db *gorm.DB) *PointCloudService {
	return &PointCloudService{DB: db}
}

func (s *PointCloudService) List(page, pageSize int) ([]models.PointCloud, int64, error) {
	var list []models.PointCloud
	var total int64

	offset := (page - 1) * pageSize

	if err := s.DB.Model(&models.PointCloud{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *PointCloudService) GetByID(id uuid.UUID) (*models.PointCloud, error) {
	var pc models.PointCloud
	if err := s.DB.First(&pc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("point cloud not found")
		}
		return nil, err
	}
	return &pc, nil
}

func (s *PointCloudService) Create(dto *CreatePointCloudDTO) (*models.PointCloud, error) {
	pc := &models.PointCloud{
		Name:       dto.Name,
		FilePath:   dto.FilePath,
		FileSize:   dto.FileSize,
		PointCount: dto.PointCount,
		BoundsMinX: dto.BoundsMinX,
		BoundsMinY: dto.BoundsMinY,
		BoundsMinZ: dto.BoundsMinZ,
		BoundsMaxX: dto.BoundsMaxX,
		BoundsMaxY: dto.BoundsMaxY,
		BoundsMaxZ: dto.BoundsMaxZ,
		UploadedBy: dto.UploadedBy,
	}

	if err := s.DB.Create(pc).Error; err != nil {
		return nil, err
	}

	return s.GetByID(pc.ID)
}

func (s *PointCloudService) Update(id uuid.UUID, dto *UpdatePointCloudDTO) (*models.PointCloud, error) {
	pc, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	pc.Name = dto.Name
	pc.FilePath = dto.FilePath
	pc.FileSize = dto.FileSize
	pc.PointCount = dto.PointCount

	if err := s.DB.Save(pc).Error; err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *PointCloudService) Delete(id uuid.UUID) error {
	result := s.DB.Delete(&models.PointCloud{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("point cloud not found")
	}
	return nil
}
