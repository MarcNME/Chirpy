package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) userHandler(w http.ResponseWriter, r *http.Request) {
	var userReq createUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Could not decode body\n" + err.Error()))
		if err != nil {
			return
		}

		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userReq.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error creating user\n" + err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error marshalling user\n" + err.Error()))
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}

type createUserRequest struct {
	Email string `json:"email"`
}
