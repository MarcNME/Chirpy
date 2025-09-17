package mappers

import (
	"github.com/MarcNME/Chirpy/internal/gen/database"
	"github.com/MarcNME/Chirpy/internal/models"
)

func UserToDTO(u database.User) models.UserDTO {
	return models.UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
