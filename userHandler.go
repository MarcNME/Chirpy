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
	var userReq createOrUpdateUserRequest
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

	userDTOJson, err := json.Marshal(mappers.UserToDTO(user))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling user\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userDTOJson)
	if err != nil {
		return
	}
}

func (cfg *apiConfig) updateUserMailAndPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Unauthorized\n"+err.Error(), http.StatusUnauthorized)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		helpers.WriteErrorMessage(w, "Unauthorized\n"+err.Error(), http.StatusUnauthorized)
		return
	}

	var userReq createOrUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		helpers.WriteErrorMessage(w, "Could not decode body\n"+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error hashing password\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := cfg.db.UpdateUserEmailAndPassword(r.Context(), userId, userReq.Email, hashedPassword)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error updating user\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	userDTOJson, err := json.Marshal(mappers.UserToDTO(updatedUser))
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

type createOrUpdateUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
