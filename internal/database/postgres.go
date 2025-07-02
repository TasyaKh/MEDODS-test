package database

import (
	"fmt"
	"test-task/internal/config"
	"test-task/internal/models"
	"time"

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

	var db *gorm.DB
	var err error

	const maxAttempts = 5
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("DB connection failed attempt %d , err: %v ", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to DB")
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
