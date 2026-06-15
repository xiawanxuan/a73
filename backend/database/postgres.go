package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ocean-survey-backend/config"
	"ocean-survey-backend/models"
)

func InitPostgres(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.PointCloud{},
		&models.TerrainLabel{},
		&models.Annotation{},
		&models.AnnotationSnapshot{},
	)
	if err != nil {
		return nil, err
	}

	if err := SeedTerrainLabels(db); err != nil {
		log.Printf("Warning: failed to seed terrain labels: %v", err)
	}

	return db, nil
}

func SeedTerrainLabels(db *gorm.DB) error {
	labels := []models.TerrainLabel{
		{Name: "海沟", Color: "#1e3a5f", Description: "海底深陷的狭长沟槽地形"},
		{Name: "礁石", Color: "#8b6914", Description: "海中突出的岩石或珊瑚礁"},
		{Name: "水下管线", Color: "#4a5568", Description: "铺设在海底的管道或线缆"},
	}

	for i := range labels {
		result := db.Where("name = ?", labels[i].Name).FirstOrCreate(&labels[i])
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
