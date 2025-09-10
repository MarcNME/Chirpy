package main

import (
	"fmt"
	"net/http"
)

const template = "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>"

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf(template, cfg.fileserverHits.Load())
	r.Header.Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(message))
	if err != nil {
		return
	}
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
