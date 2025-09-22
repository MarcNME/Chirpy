package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/MarcNME/Chirpy/internal/auth"
	"github.com/MarcNME/Chirpy/internal/mappers"
	"github.com/google/uuid"
)

var profanityWords = [...]string{"kerfuffle", "sharbert", "fornax"}

const ErrorWritingResponse = "Error writing response: %v"

func (cfg *apiConfig) NewChirpHandler(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Missing or invalid Authorization header: "+err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirpRequest := newChirpRequest{}
	err = decoder.Decode(&chirpRequest)
	if err != nil {
		helpers.WriteErrorMessage(w, "Could not decode body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(chirpRequest.Body) > 140 {
		helpers.WriteErrorMessage(w, "Chirp is to long", http.StatusBadRequest)
		return
	}

	msg := helpers.ReplaceAllIgnoreCase(chirpRequest.Body, profanityWords[:])
	chirp, err := cfg.db.CreateChirp(r.Context(), msg, userId)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error creating chirp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	chirpDtoJson, err := json.Marshal(mappers.ChirpToDTO(chirp))
	if err != nil {
		log.Printf("Error marshalling chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(chirpDtoJson)
	if err != nil {
		log.Printf(ErrorWritingResponse, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		helpers.WriteErrorMessage(w, "Error getting chirps: "+err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(mappers.ChirpsToDTOs(chirps))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling chirps: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf(ErrorWritingResponse, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) GetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("id")
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error getting chirp: "+err.Error(), http.StatusNotFound)
		return
	}

	body, err := json.Marshal(mappers.ChirpToDTO(chirp))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling chirp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf(ErrorWritingResponse, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) DeleteChirpById(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Missing or invalid Authorization header: "+err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	chirpId := r.PathValue("id")
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		helpers.WriteErrorMessage(w, fmt.Sprintf("Can not find Chirp with id %v: %v", chirpId, err.Error()), http.StatusNotFound)
		return
	}

	if chirp.UserID != userId {
		helpers.WriteErrorMessage(w, "User is not allowed to remove Chirp", http.StatusForbidden)
		return
	}

	if err := cfg.db.DeleteChirpByID(r.Context(), chirp.ID); err != nil {
		helpers.WriteErrorMessage(w, "Error deleting Chirp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type newChirpRequest struct {
	Body string `json:"body"`
}
