package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/helpers"
)

var profanityWords = [...]string{"kerfuffle", "sharbert", "fornax"}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	validationRequest := validationRequest{}
	err := decoder.Decode(&validationRequest)

	if err != nil {
		log.Printf("Error decoding body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		body, err := json.Marshal(errorResponse{Error: err.Error()})
		if err != nil {
			log.Printf("Error marshalling error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(body), http.StatusBadRequest)
		return
	}

	if len(validationRequest.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		body, err := json.Marshal(errorResponse{Error: "Chirp is to long"})
		if err != nil {
			log.Printf("Error marshalling error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			log.Printf("Error writing response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	msg := helpers.ReplaceAllIgnoreCase(validationRequest.Body, profanityWords[:])
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"cleaned_body": "` + msg + `"}`))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type validationRequest struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}
