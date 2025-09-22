package helpers

import (
	"log"
	"net/http"

	"github.com/MarcNME/Chirpy/constants"
)

func WriteErrorMessage(w http.ResponseWriter, msg string, errorCode int) {
	log.Printf("%s: %d", msg, errorCode)
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(errorCode)
	_, err := w.Write([]byte(`{"error": "` + msg + `"}`))
	if err != nil {
		log.Printf("Could not write error message: %v", err)
		return
	}
}
