package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeaderValue := headers.Get("Authorization")

	if authHeaderValue == "" {
		return "", fmt.Errorf("expected header is empty")
	}
	if !strings.HasPrefix(authHeaderValue, "ApiKey ") {
		return "", fmt.Errorf("invalid authorization header")
	}
	return strings.TrimPrefix(authHeaderValue, "ApiKey "), nil
}
