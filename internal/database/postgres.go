package database

import (
	"fmt"
	"test-task/internal/config"
	"test-task/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if cfg.AutoMigrate {
		err = db.AutoMigrate(&models.User{}, &models.Session{})
		if err != nil {
			return nil, err
		}
		fmt.Println("AutoMigrate completed")
	}
	return db, nil
}
