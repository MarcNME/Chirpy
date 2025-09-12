package main

import (
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
)

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(constants.ContentType, constants.TextPlain)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
