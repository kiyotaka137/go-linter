package rules

import "strings"

func FindSensitiveKeyword(msg string, keywords []string) (string, bool) {
	msgLower := strings.ToLower(msg)

	for _, kw := range keywords {
		kw = strings.TrimSpace(strings.ToLower(kw))
		if kw == "" {
			continue
		}

		if strings.Contains(msgLower, kw) {
			return kw, true
		}
	}

	return "", false
}
