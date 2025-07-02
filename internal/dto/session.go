package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionCreate struct {
	UserID           uuid.UUID
	JTI              string
	RefreshTokenHash string
	UserAgent        string
	IP               string
	ExpiresAt        time.Time
}
