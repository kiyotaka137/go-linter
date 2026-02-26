package loglint

import (
	"fmt"
	"go-linter/analyzer"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct {
	Config string `json:"config"`
}

type Plugin struct {
	settings Settings
}

func New(settings any) (register.LinterPlugin, error) {
	cfg, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: cfg}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.Config != "" {
		if err := analyzer.Analyzer.Flags.Set("config", p.settings.Config); err != nil {
			return nil, fmt.Errorf("set analyzer config: %w", err)
		}
	}

	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
