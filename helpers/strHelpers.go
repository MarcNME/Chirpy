package helpers

import (
	"strings"
)

func ContainsAny(str string, substrs []string) bool {
	for _, s := range substrs {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func ReplaceAllIgnoreCase(str string, substrs []string) string {
	strList := strings.Split(str, " ")
	var out string
	for _, s := range strList {
		if out != "" {
			out += " "
		}
		found := false
		for _, word := range substrs {
			if strings.ToLower(s) == strings.ToLower(word) {
				found = true
			}
		}
		if found {
			out += "****"
		} else {
			out += s
		}
	}
	return out
}
