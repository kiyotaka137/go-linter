package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Rules             RulesConfig `yaml:"rules"`
	SensitiveKeywords []string    `yaml:"sensitive_keywords"`
}

type RulesConfig struct {
	Lowercase RuleConfig          `yaml:"lowercase"`
	Symbols   RuleConfig          `yaml:"symbols"`
	Sensitive SensitiveRuleConfig `yaml:"sensitive"`
}

type RuleConfig struct {
	Enabled *bool `yaml:"enabled"`
}

type SensitiveRuleConfig struct {
	Enabled  *bool    `yaml:"enabled"`
	Keywords []string `yaml:"keywords"`
	Patterns []string `yaml:"patterns"`
}

type SensitivePattern struct {
	Source string
	Regex  *regexp.Regexp
}

type RuntimeConfig struct {
	EnableLowercase   bool
	EnableSymbols     bool
	EnableSensitive   bool
	SensitiveWords    []string
	SensitivePatterns []SensitivePattern
}

var DefaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"token",
	"secret",
	"api_key",
	"apikey",
	"jwt",
	"bearer",
	"authorization",
	"cookie",
	"session",
	"private_key",
	"client_secret",
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		EnableLowercase: true,
		EnableSymbols:   true,
		EnableSensitive: true,
		SensitiveWords:  append([]string(nil), DefaultSensitiveKeywords...),
	}
}

func Load(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	cfg.SensitiveKeywords = NormalizeKeywords(cfg.SensitiveKeywords)

	return &cfg, nil
}

func Resolve(cfg *Config) RuntimeConfig {
	resolved := DefaultRuntimeConfig()
	if cfg == nil {
		resolved.SensitiveWords = NormalizeKeywords(resolved.SensitiveWords)
		return resolved
	}

	resolved.EnableLowercase = boolOrDefault(cfg.Rules.Lowercase.Enabled, true)
	resolved.EnableSymbols = boolOrDefault(cfg.Rules.Symbols.Enabled, true)
	resolved.EnableSensitive = boolOrDefault(cfg.Rules.Sensitive.Enabled, true)

	switch {
	case len(cfg.Rules.Sensitive.Keywords) > 0:
		resolved.SensitiveWords = NormalizeKeywords(cfg.Rules.Sensitive.Keywords)
	case len(cfg.SensitiveKeywords) > 0:
		resolved.SensitiveWords = NormalizeKeywords(cfg.SensitiveKeywords)
	default:
		resolved.SensitiveWords = NormalizeKeywords(resolved.SensitiveWords)
	}
	resolved.SensitivePatterns = CompilePatterns(cfg.Rules.Sensitive.Patterns)

	return resolved
}

func boolOrDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}

	return *value
}

func NormalizeKeywords(keywords []string) []string {
	seen := make(map[string]struct{}, len(keywords))
	normalized := make([]string, 0, len(keywords))

	for _, kw := range keywords {
		kw = strings.TrimSpace(strings.ToLower(kw))
		if kw == "" {
			continue
		}

		if _, ok := seen[kw]; ok {
			continue
		}

		seen[kw] = struct{}{}
		normalized = append(normalized, kw)
	}

	return normalized
}

func CompilePatterns(patterns []string) []SensitivePattern {
	seen := make(map[string]struct{}, len(patterns))
	compiled := make([]SensitivePattern, 0, len(patterns))

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if _, ok := seen[pattern]; ok {
			continue
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		seen[pattern] = struct{}{}
		compiled = append(compiled, SensitivePattern{
			Source: pattern,
			Regex:  re,
		})
	}

	return compiled
}
