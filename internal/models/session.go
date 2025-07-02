package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uint      `gorm:"primaryKey"`
	UserID           uuid.UUID `gorm:"not null"`
	JTI              string    `gorm:"unique;not null"` // JWT ID
	RefreshTokenHash string    `gorm:"not null"`
	UserAgent        string    `gorm:"not null"`
	IP               string    `gorm:"not null"`
	CreatedAt        time.Time
	RefreshExpiresIn time.Time
}
