package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var id, _ = uuid.Parse("47734160-3C16-4E82-AA42-7EEC620B2802")
var signature = "DA3B2880-1519-45C9-81BB-ABCA879DEC65"

func TestMakeJWT(t *testing.T) {
	_, err := MakeJWT(id, signature, time.Minute*15)
	if err != nil {
		t.Errorf("MakeJWT() returned error \"%v\"; want nil", err)
		return
	}
}

func TestValidateJWT(t *testing.T) {
	testData := []struct {
		token string
		error error
	}{
		{
			// Unexpired token update before 20.11.2286 - 18:46:39 UTC
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI0NzczNDE2MC0zYzE2LTRlODItYWE0Mi03ZWVjNjIwYjI4MDIiLCJleHAiOjk5OTk5OTk5OTksImlhdCI6MTc1ODE4NTM2Nn0.YcENIaO8UbX3vXMkNLHROoXPwSSyKyt70GHClbISBl8",
			error: nil,
		},
		{
			// Expired token
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI0NzczNDE2MC0zYzE2LTRlODItYWE0Mi03ZWVjNjIwYjI4MDIiLCJleHAiOjE3NTgxODUzNjYsImlhdCI6MTc1ODE4NTM2Nn0.5gl_B79AmBM9TwmKu0SbwcKLq0FK10Memx4nqAJ1SIc",
			error: fmt.Errorf("token has invalid claims: token is expired"),
		},
		{
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI0NzczNDE2MC0zYzE2LTRlODItYWE0Mi03ZWVjNjIwYjI4MDIiLCJleHAiOjk5OTk5OTk5OTksImlhdCI6MTc1ODE4NTM2Nn0.wzoqrWK7VGsur21j02u3Vu1fSAeoGIlQ-azTHWS4jkA",
			error: fmt.Errorf("token signature is invalid: signature is invalid"),
		},
	}

	for _, test := range testData {
		resultId, err := ValidateJWT(test.token, signature)
		if test.error == nil {
			if err != nil {
				t.Errorf("ValidateJWT() returned error \"%v\"; want \"nil\"", err)
				return
			}
			if resultId != id {
				t.Errorf("ValidateJWT() returned id %v; want %v", resultId, id)
			}
		} else {
			if err.Error() != test.error.Error() {
				t.Errorf("ValidateJWT() returned error \"%v\"; want \"%v\"", err, test.error)
				return
			}
		}
	}
}
