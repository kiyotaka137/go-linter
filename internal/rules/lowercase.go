package rules

import (
	"strings"
	"unicode"
)

func IsLowercase(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}

	firstCh := []rune(trimmed)[0]
	return unicode.IsLower(firstCh)
}
