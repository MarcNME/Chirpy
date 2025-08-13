package handlers

import (
	"net/http"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
