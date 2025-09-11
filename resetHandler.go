package main

import "net/http"

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "DEV" {
		err := cfg.db.DeleteAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error deleting users\n" + err.Error()))
			return
		}

		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden"))
	}
}
