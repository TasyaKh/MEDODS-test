package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
}
