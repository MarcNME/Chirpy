package mappers

import (
	"github.com/MarcNME/Chirpy/internal/gen/database"
	"github.com/MarcNME/Chirpy/internal/models"
)

func ChirpToDTO(c database.Chirp) models.ChirpDTO {
	return models.ChirpDTO{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

func ChirpsToDTOs(c []database.Chirp) []models.ChirpDTO {
	var dtos = make([]models.ChirpDTO, 0, len(c))

	for _, chirp := range c {
		dtos = append(dtos, ChirpToDTO(chirp))
	}

	return dtos
}
