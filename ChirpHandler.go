package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/google/uuid"
)

var profanityWords = [...]string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) NewChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	validationRequest := validationRequest{}
	err := decoder.Decode(&validationRequest)
	if err != nil {
		helpers.WriteErrorMessage(w, "Could not decode body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(validationRequest.Body) > 140 {
		helpers.WriteErrorMessage(w, "Chirp is to long", http.StatusBadRequest)
		return
	}

	msg := helpers.ReplaceAllIgnoreCase(validationRequest.Body, profanityWords[:])
	chirp, err := cfg.db.CreateChirp(r.Context(), msg, validationRequest.UserID)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error creating chirp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(res)
	if err != nil {
		log.Printf("Error writing response: %v", err)
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

	body, err := json.Marshal(chirps)
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
	chirpId := r.PathValue("id")
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		helpers.WriteErrorMessage(w, "Error getting chirp: "+err.Error(), http.StatusNotFound)
		return
	}

	body, err := json.Marshal(chirp)
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

type validationRequest struct {
	Body   string        `json:"body"`
	UserID uuid.NullUUID `json:"user_id"`
}
