package main

import (
	"encoding/json"
	"net/http"

	"github.com/MarcNME/Chirpy/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) UpgradeUserToChirpyRed(w http.ResponseWriter, r *http.Request) {
	var request updateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "could not decode body"+err.Error(), http.StatusBadRequest)
	}

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err := cfg.db.UpdateUserToChirpyRed(r.Context(), request.Data.UserId)
	if err != nil {
		helpers.WriteErrorMessage(w, "could not update user: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type updateSubscriptionRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}
