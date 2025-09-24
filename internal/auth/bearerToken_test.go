package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	testData := []struct {
		headers  http.Header
		expected string
		err      error
	}{
		{
			headers:  http.Header{"Authorization": {"Bearer mytoken", "othertoken"}},
			expected: "mytoken",
			err:      nil,
		},
		{
			headers:  http.Header{"Authorization": {"Bearer mytoken"}},
			expected: "mytoken",
			err:      nil,
		},
		{
			headers:  http.Header{"Authorization": {"mytoken"}},
			expected: "mytoken",
			err:      fmt.Errorf("invalid authorization header"),
		},
		{
			headers:  http.Header{},
			expected: "",
			err:      fmt.Errorf("expected header is empty"),
		},
	}

	for _, test := range testData {
		result, err := GetBearerToken(test.headers)
		if test.err == nil {
			if err != nil {
				t.Errorf("GetBearerToken(%v) returned error \"%v\"; want nil", test.headers, err)
				return
			}

			if result != test.expected {
				t.Errorf("GetBearerToken(%v) = %q; want %q", test.headers, result, test.expected)
			}
		} else {
			if err.Error() != test.err.Error() {
				t.Errorf("GetBearerToken(%v) returned error \"%v\"; want \"%v\"", test.headers, err, test.err)
				return
			}
		}
	}
}
