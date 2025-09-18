package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeaderValue := headers.Get("Authorization")

	if authHeaderValue == "" {
		return "", fmt.Errorf("expected header is empty")
	}

	return strings.TrimPrefix(authHeaderValue, "Bearer "), nil
}
