package auth

import "testing"

func TestMakeRefreshToken(t *testing.T) {
	for i := 0; i < 10; i++ {
		token, err := MakeRefreshToken()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(token) != 64 {
			t.Fatalf("Expected length 32, got %v", len(token))
		}
	}
}
