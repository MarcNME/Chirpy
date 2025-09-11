package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MarcNME/Chirpy/internal/database"
	"github.com/joho/godotenv"
)
import _ "github.com/lib/pq"

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	dbQueries := database.New(db)

	fmt.Println(dbQueries)
	cfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(FileHandler())))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.metricsResetHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: Log(mux),
		HTTP2:   &http.HTTP2Config{},
	}

	log.Printf("Serving on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
