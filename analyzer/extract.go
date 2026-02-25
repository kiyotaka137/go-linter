package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

var logMethods = map[string]struct{}{
	"Debug": {},
	"Info":  {},
	"Warn":  {},
	"Error": {},
}

func extractLogMessageArg(call *ast.CallExpr) (ast.Expr, bool) {
	if !isLogCall(call) || len(call.Args) == 0 {
		return nil, false
	}

	return call.Args[0], true
}

func isLogCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if _, ok := logMethods[selector.Sel.Name]; !ok {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "log" || ident.Name == "slog"
}

func extractStaticString(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}

		value, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", false
		}

		return value, true
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}

		left, ok := extractStaticString(e.X)
		if !ok {
			return "", false
		}

		right, ok := extractStaticString(e.Y)
		if !ok {
			return "", false
		}

		return left + right, true
	case *ast.ParenExpr:
		return extractStaticString(e.X)
	default:
		return "", false
	}
}
