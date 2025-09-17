package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	passwordHashes := []struct {
		password string
		error    error
	}{
		{
			password: "123456",
			error:    nil,
		},
		{
			password: "password",
			error:    nil,
		},
	}

	for _, test := range passwordHashes {
		result, err := HashPassword(test.password)
		if err != nil {
			t.Errorf("HashPassword(%q) returned error %v; want nil", test.password, err)
			return
		}
		err = CheckPasswordHash(test.password, result)
		if err != nil {
			t.Errorf("CheckPasswordHash(%q, %q) returned error %v; want nil", test.password, result, err)
		}
	}
}
