package main

import (
	"encoding/json"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/MarcNME/Chirpy/internal/auth"
	"github.com/MarcNME/Chirpy/internal/mappers"
)

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	var userReq createUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		helpers.WriteErrorMessage(w, "Could not decode body\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if userReq.Email == "" || userReq.Password == "" {
		helpers.WriteErrorMessage(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error hashing password\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userReq.Email, hashedPassword)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error creating user\n"+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	userDTOJson, err := json.Marshal(mappers.UserToDTO(user))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling user\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	_, err = w.Write(userDTOJson)
	if err != nil {
		return
	}
}

type createUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
