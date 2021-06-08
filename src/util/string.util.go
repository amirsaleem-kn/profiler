package util

import "strings"

// to camelCase
func ToCamelCase(s string) string {
	var len = len(s)

	// TODO: instead of this check, add a more robust check
	// check if string contains all lowercase characters except the first character
	// if yes, turn it into a lowercase string
	if len == 2 {
		s = strings.ToLower(s)
		return s
	}

	s = strings.ToLower(s[:1]) + s[1:]
	return s
}
