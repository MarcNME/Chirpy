package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currentDate := &jwt.NumericDate{
		Time: time.Now().UTC(),
	}
	expireDate := &jwt.NumericDate{
		Time: time.Now().Add(expiresIn).UTC(),
	}
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userId.String(),
		ExpiresAt: expireDate,
		IssuedAt:  currentDate,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return [16]byte{}, err
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return [16]byte{}, err
	}

	return id, nil
}
