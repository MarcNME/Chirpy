package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/helpers"
	"github.com/google/uuid"
)

var profanityWords = [...]string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) ChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	validationRequest := validationRequest{}
	err := decoder.Decode(&validationRequest)
	if err != nil {
		writeErrorMessage(w, "Could not decode body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(validationRequest.Body) > 140 {
		writeErrorMessage(w, "Chirp is to long", http.StatusBadRequest)
		return
	}

	msg := helpers.ReplaceAllIgnoreCase(validationRequest.Body, profanityWords[:])
	chirp, err := cfg.db.CreateChirp(r.Context(), msg, validationRequest.UserID)
	if err != nil {
		writeErrorMessage(w, "Error creating chirp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeErrorMessage(w http.ResponseWriter, msg string, errorCode int) {
	w.WriteHeader(errorCode)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"error": "` + msg + `"}`))
	if err != nil {
		log.Printf("Could not write error message: %v", err)
		return
	}
}

type validationRequest struct {
	Body   string        `json:"body"`
	UserID uuid.NullUUID `json:"user_id"`
}
