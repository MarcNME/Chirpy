package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID    `json:"id"`
	Email     string       `json:"email"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
