package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MarcNME/Chirpy/internal/gen/database"
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
	cfg := getApiConfig()

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(FileHandler())))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.metricsResetHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	// User endpoints
	mux.HandleFunc("POST /api/users", cfg.userHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateUserMailAndPasswordHandler)
	// login endpoints
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.refreshAuthTokenHandler)
	mux.HandleFunc("POST /api/revoke", cfg.revokeRefreshToken)
	// Chirp endpoints
	mux.HandleFunc("POST /api/chirps", cfg.NewChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.DeleteChirpById)

	srv := &http.Server{
		Addr:    cfg.address + ":" + cfg.port,
		Handler: Log(mux),
		HTTP2:   &http.HTTP2Config{},
	}

	log.Printf("Serving on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getApiConfig() *apiConfig {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		log.Println("Only using environment variables")
	}

	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	dbQueries := database.New(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		address:        addr,
		port:           port,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      jwtSecret,
	}

	return cfg
}

type apiConfig struct {
	fileserverHits atomic.Int32
	address        string
	port           string
	db             *database.Queries
	platform       string
	jwtSecret      string
}
