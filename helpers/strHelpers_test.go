package helpers

import (
	"testing"
)

func TestContainsAny(t *testing.T) {
	tests := []struct {
		str      string
		substrs  []string
		expected bool
	}{
		{"hello world", []string{"world", "universe"}, true},
		{"hello world", []string{"universe", "galaxy"}, false},
		{"golang is great", []string{"go", "lang"}, true},
		{"golang is great", []string{"python", "java"}, false},
		{"test string", []string{}, false},
	}

	for _, test := range tests {
		result := ContainsAny(test.str, test.substrs)
		if result != test.expected {
			t.Errorf("ContainsAny(%q, %v) = %v; want %v", test.str, test.substrs, result, test.expected)
		}
	}
}

func TestReplaceAllIgnoreCase(t *testing.T) {
	tests := []struct {
		input    string
		substrs  []string
		expected string
	}{
		{"hello world", []string{"world"}, "hello ****"},
		{"Hello World", []string{"world"}, "Hello ****"},
		{"This is a kerfuffle opinion I need to share with the world", []string{"kerfuffle"}, "This is a **** opinion I need to share with the world"},
		{"This is a kerfuffle! opinion I need to share with the world", []string{"kerfuffle"}, "This is a kerfuffle! opinion I need to share with the world"},
		{"Hello there general Kenobi", []string{"general", "Kenobi"}, "Hello there **** ****"},
	}

	for _, test := range tests {
		result := ReplaceAllIgnoreCase(test.input, test.substrs)
		if result != test.expected {
			t.Errorf("ReplaceAll(%q, %v) = %q; want %q", test.input, test.substrs, result, test.expected)
		}
	}
}
