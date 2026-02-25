package analyzer

import (
	"go-linter/internal/config"
	"go-linter/internal/rules"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "reports logs",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	keywords := config.DefaultSensitiveKeywords
	if cfg, err := config.Load("config.yaml"); err == nil && len(cfg.SensitiveKeywords) > 0 {
		keywords = cfg.SensitiveKeywords
	}
	keywords = config.NormalizeKeywords(keywords)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			msgExpr, ok := extractLogMessageArg(call)
			if !ok {
				return true
			}

			msg, ok := extractStaticString(msgExpr)
			if !ok {
				return true
			}

			if keyword, found := rules.FindSensitiveKeyword(msg, keywords); found {
				pass.Reportf(msgExpr.Pos(), "log message have sensitive keyword %q", keyword)
			}

			if !rules.IsLowercase(msg) {
				pass.Reportf(msgExpr.Pos(), "log message must start with lowercase letter")
			}
			if !rules.IsWithoutSymbols(msg) {
				pass.Reportf(msgExpr.Pos(), "log message have forbidden symbols")
			}

			return true
		})
	}

	return nil, nil
}
