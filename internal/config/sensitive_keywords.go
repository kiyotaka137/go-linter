package config

import (
	"fmt"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SensitiveKeywords []string `yaml:"sensitive_keywords"`
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

func Load(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	cfg.SensitiveKeywords = NormalizeKeywords(cfg.SensitiveKeywords)

	return &cfg, nil
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
