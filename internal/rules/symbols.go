package rules

import "strings"

func isEnglishCh(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsWithoutSymbols(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}

	for _, r := range trimmed {
		if !isNumber(r) && !isEnglishCh(r) && r != ' ' {
			return false
		}
	}

	return true
}
