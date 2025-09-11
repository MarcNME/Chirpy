package models

import (
	"time"

	"github.com/MarcNME/Chirpy/internal/database"
	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func ToUserDTO(u database.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email,
	}
}
