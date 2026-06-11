package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the parsed config.yaml. The known keys are pulled out for the
// generator's own use; Site carries every key (known and extra) for exposure
// to templates as site.*.
type Config struct {
	SiteName  string
	BaseURL   string
	OutputDir string
	Site      map[string]any
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if raw == nil {
		raw = map[string]any{}
	}

	cfg := &Config{
		OutputDir: "site",
		Site:      raw,
	}
	if v, ok := raw["site_name"].(string); ok {
		cfg.SiteName = v
	}
	if v, ok := raw["base_url"].(string); ok {
		cfg.BaseURL = v
	}
	if v, ok := raw["output_dir"].(string); ok && v != "" {
		cfg.OutputDir = v
	}
	return cfg, nil
}
