package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/MarcNME/Chirpy/internal/auth"
	"github.com/MarcNME/Chirpy/internal/mappers"
	"github.com/google/uuid"
)

var profanityWords = [...]string{"kerfuffle", "sharbert", "fornax"}

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
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Missing or invalid Authorization header: "+err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

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
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) GetChirpById(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Missing or invalid Authorization header: "+err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

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
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type newChirpRequest struct {
	Body string `json:"body"`
}
