package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/MarcNME/Chirpy/internal/auth"
	"github.com/MarcNME/Chirpy/internal/mappers"
	"github.com/MarcNME/Chirpy/internal/models"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if req.ExpiresInSeconds <= 0 || req.ExpiresInSeconds > 3600 {
		req.ExpiresInSeconds = 3600 // Default to 1 hour
	}

	duration := time.Duration(req.ExpiresInSeconds) * time.Second

	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		helpers.WriteErrorMessage(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err = auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		helpers.WriteErrorMessage(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, duration)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := loginResponse{
		UserDTO: mappers.UserToDTO(user),
		Token:   token,
	}

	userDTOJson, err := json.Marshal(response)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling user\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userDTOJson)
	if err != nil {
		return
	}
}

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int64  `json:"expires_in_seconds,omitempty"`
}

type loginResponse struct {
	models.UserDTO
	Token string `json:"token"`
}
