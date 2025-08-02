package database

import (
	"fmt"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/vert3xc/bebrochka/server/internal/models"
	"github.com/vert3xc/bebrochka/server/config"
)

var DB *gorm.DB

func InitDatabase(cfg config.Config) (*gorm.DB, error) {
	db, err := NewPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("DB error: %v", err)
	}

	err = db.AutoMigrate(&models.Flag{})
	if err != nil {
		return nil, fmt.Errorf("AutoMigrate failed: %v", err)
	}
	return db, nil
}