package analyzer

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	if err := Analyzer.Flags.Set("config", filepath.Join("..", "config.yaml")); err != nil {
		t.Fatalf("set config flag: %v", err)
	}

	dir := filepath.Join("..", "testdata")
	analysistest.Run(t, dir, Analyzer,
		"./src/badlogs_logslog",
		"./src/badlogs_zap",
	)
}
