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

var configPath = "config.yaml"

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", configPath, "path to YAML config file")
}

func run(pass *analysis.Pass) (interface{}, error) {
	runtimeCfg := config.DefaultRuntimeConfig()
	if cfg, err := config.Load(configPath); err == nil {
		runtimeCfg = config.Resolve(cfg)
	}

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

			if runtimeCfg.EnableSensitive {
				if keyword, found := rules.FindSensitiveKeyword(msg, runtimeCfg.SensitiveWords); found {
					pass.Reportf(msgExpr.Pos(), "log message have sensitive keyword %q", keyword)
				}
				if pattern, found := findSensitivePattern(msg, runtimeCfg.SensitivePatterns); found {
					pass.Reportf(msgExpr.Pos(), "log message matches sensitive pattern %q", pattern)
				}
			}
			if runtimeCfg.EnableLowercase && !rules.IsLowercase(msg) {
				fixed, ok := suggestLowercaseFix(msg)
				reportWithFix(pass, msgExpr, "log message must start with lowercase letter", fixed, ok)
			}
			if runtimeCfg.EnableSymbols && !rules.IsWithoutSymbols(msg) {
				fixed, ok := suggestSymbolsFix(msg)
				reportWithFix(pass, msgExpr, "log message have forbidden symbols", fixed, ok)
			}

			return true
		})
	}

	return nil, nil
}
