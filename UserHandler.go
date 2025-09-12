package main

import (
	"encoding/json"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
)

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	var userReq createUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		helpers.WriteErrorMessage(w, "Could not decode body\n"+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userReq.Email)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error creating user\n"+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(user)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling user\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}

type createUserRequest struct {
	Email string `json:"email"`
}
