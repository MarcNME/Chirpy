package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {
	testData := []struct {
		headers  http.Header
		expected string
		err      error
	}{
		{
			headers:  http.Header{"Authorization": {"ApiKey myKey", "myOtherKey"}},
			expected: "myKey",
			err:      nil,
		},
		{
			headers:  http.Header{"Authorization": {"ApiKey myKey"}},
			expected: "myKey",
			err:      nil,
		},
		{
			headers:  http.Header{"Authorization": {"myKey"}},
			expected: "",
			err:      fmt.Errorf("invalid authorization header"),
		},
		{
			headers:  http.Header{},
			expected: "",
			err:      fmt.Errorf("expected header is empty"),
		},
	}

	for _, test := range testData {
		result, err := GetApiKey(test.headers)
		if test.err == nil {
			if err != nil {
				t.Errorf("GetApiKey(%v) returned error \"%v\"; want nil", test.headers, err)
				return
			}

			if result != test.expected {
				t.Errorf("GetApiKey(%v) = %q; want %q", test.headers, result, test.expected)
			}
		} else {
			if err.Error() != test.err.Error() {
				t.Errorf("GetApiKey(%v) returned error \"%v\"; want \"%v\"", test.headers, err, test.err)
				return
			}
		}
	}
}
