package main

import (
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/handlers"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./static/"))))
	mux.HandleFunc("/healthz", handlers.ReadinessHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: Log(mux),
		HTTP2:   &http.HTTP2Config{},
	}

	log.Printf("Serving on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
