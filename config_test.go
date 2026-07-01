package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadConfig(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, `
site_name: Acme
base_url: https://example.com
output_dir: dist
phone: "555-1234"
nav:
  - Home
  - About
`))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SiteName != "Acme" {
		t.Errorf("SiteName = %q", cfg.SiteName)
	}
	if cfg.BaseURL != "https://example.com" {
		t.Errorf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.OutputDir != "dist" {
		t.Errorf("OutputDir = %q", cfg.OutputDir)
	}
	if cfg.Site["phone"] != "555-1234" {
		t.Errorf("Site[phone] = %v", cfg.Site["phone"])
	}
	if cfg.Site["site_name"] != "Acme" {
		t.Errorf("known keys should also appear under site.*, got %v", cfg.Site["site_name"])
	}
	nav, ok := cfg.Site["nav"].([]any)
	if !ok || len(nav) != 2 {
		t.Errorf("Site[nav] = %v", cfg.Site["nav"])
	}
}

func TestLoadConfigNormalizesBaseURL(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, "base_url: https://example.com/\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://example.com" {
		t.Errorf("BaseURL = %q, want trailing slash stripped", cfg.BaseURL)
	}
	if cfg.Site["base_url"] != "https://example.com" {
		t.Errorf("Site[base_url] = %v, want normalized value for templates", cfg.Site["base_url"])
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, "site_name: X\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.OutputDir != "site" {
		t.Errorf("OutputDir default = %q, want site", cfg.OutputDir)
	}
}

func TestLoadConfigEmpty(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Site == nil {
		t.Error("Site map should not be nil for empty config")
	}
	if cfg.OutputDir != "site" {
		t.Errorf("OutputDir default = %q, want site", cfg.OutputDir)
	}
}

func TestLoadConfigMissing(t *testing.T) {
	if _, err := loadConfig(filepath.Join(t.TempDir(), "config.yaml")); err == nil {
		t.Error("expected error for missing config file")
	}
}
