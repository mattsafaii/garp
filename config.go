package main

import (
	"fmt"
	"os"
	"strings"

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
	normalizeDates(raw)

	cfg := &Config{
		OutputDir: "site",
		Site:      raw,
	}
	if v, ok := raw["site_name"].(string); ok {
		cfg.SiteName = v
	}
	// A trailing slash on base_url would double up everywhere it's joined
	// with a page URL (canonical, OG, sitemap, robots) — normalize it away,
	// in the site.* map too so templates see the same value.
	if v, ok := raw["base_url"].(string); ok {
		v = strings.TrimRight(v, "/")
		cfg.BaseURL = v
		raw["base_url"] = v
	}
	if v, ok := raw["output_dir"].(string); ok && v != "" {
		cfg.OutputDir = v
	}
	return cfg, nil
}
