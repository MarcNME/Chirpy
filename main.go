package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"
	cfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(FileHandler())))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.metricsResetHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: Log(mux),
		HTTP2:   &http.HTTP2Config{},
	}

	log.Printf("Serving on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
