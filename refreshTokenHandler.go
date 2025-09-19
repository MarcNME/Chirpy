package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MarcNME/Chirpy/constants"
	"github.com/MarcNME/Chirpy/helpers"
	"github.com/MarcNME/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshTokenStr)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		helpers.WriteErrorMessage(w, "Refresh token is expired or revoked", http.StatusUnauthorized)
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := refreshTokenResponse{
		Token: newToken,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		helpers.WriteErrorMessage(w, "Error marshalling user\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJson)
	if err != nil {
		return
	}
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.WriteErrorMessage(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), refreshTokenStr); err != nil {
		helpers.WriteErrorMessage(w, "Could not revoke refresh token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type refreshTokenResponse struct {
	Token string `json:"token"`
}
