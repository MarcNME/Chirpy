package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type ChirpDTO struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}
