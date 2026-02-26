package analyzer

import (
	"go-linter/internal/config"
	"go/ast"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func reportWithFix(pass *analysis.Pass, expr ast.Expr, message string, fixed string, ok bool) {
	if !ok {
		pass.Report(analysis.Diagnostic{
			Pos:     expr.Pos(),
			End:     expr.End(),
			Message: message,
		})
		return
	}

	if _, isBasicString := expr.(*ast.BasicLit); !isBasicString {
		pass.Report(analysis.Diagnostic{
			Pos:     expr.Pos(),
			End:     expr.End(),
			Message: message,
		})
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		End:     expr.End(),
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "replace log message",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(strconv.Quote(fixed)),
					},
				},
			},
		},
	})
}

func suggestLowercaseFix(msg string) (string, bool) {
	runes := []rune(msg)
	for i, r := range runes {
		if unicode.IsSpace(r) {
			continue
		}

		lowered := unicode.ToLower(r)
		if lowered == r {
			return "", false
		}

		runes[i] = lowered
		return string(runes), true
	}

	return "", false
}

func suggestSymbolsFix(msg string) (string, bool) {
	trimmed := strings.TrimSpace(msg)
	if trimmed == "" {
		return "", false
	}

	var builder strings.Builder
	prevSpace := false

	for _, r := range trimmed {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
			prevSpace = false
		case r >= 'A' && r <= 'Z':
			builder.WriteRune(unicode.ToLower(r))
			prevSpace = false
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
			prevSpace = false
		case unicode.IsSpace(r):
			if builder.Len() > 0 && !prevSpace {
				builder.WriteByte(' ')
				prevSpace = true
			}
		default:
			if builder.Len() > 0 && !prevSpace {
				builder.WriteByte(' ')
				prevSpace = true
			}
		}
	}

	fixed := strings.TrimSpace(builder.String())
	if fixed == "" || fixed == msg {
		return "", false
	}

	return fixed, true
}

func findSensitivePattern(msg string, patterns []config.SensitivePattern) (string, bool) {
	for _, pattern := range patterns {
		if pattern.Regex == nil {
			continue
		}
		if pattern.Regex.MatchString(msg) {
			return pattern.Source, true
		}
	}

	return "", false
}
