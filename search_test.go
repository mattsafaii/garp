package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSearchIndexRequiresConfig(t *testing.T) {
	if _, err := runSearchIndex(t.TempDir()); err == nil {
		t.Error("expected an error without config.yaml")
	}
}

func TestRunSearchIndexRequiresBuiltSite(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")

	_, err := runSearchIndex(root)
	if err == nil {
		t.Fatal("expected an error when site/ hasn't been built yet")
	}
	if !strings.Contains(err.Error(), "run `garp build` first") {
		t.Errorf("error should mention building first: %v", err)
	}
}

func TestRunSearchIndexRequiresPagefindBinary(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	if err := os.MkdirAll(filepath.Join(root, "site"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", t.TempDir()) // a PATH with no pagefind on it, regardless of the host

	_, err := runSearchIndex(root)
	if err == nil {
		t.Fatal("expected an error when pagefind isn't installed")
	}
	if !strings.Contains(err.Error(), "pagefind not found") {
		t.Errorf("error should mention pagefind: %v", err)
	}
}

// TestRunSearchIndexWritesToStatic runs the real pagefind binary end to end
// (skipped where it isn't installed) and confirms the index lands in
// static/pagefind/ — not site/pagefind/ — so a plain `garp build` reproduces
// it in CI without pagefind installed there (the CI-purity line).
func TestRunSearchIndexWritesToStatic(t *testing.T) {
	if _, err := exec.LookPath("pagefind"); err != nil {
		t.Skip("pagefind not installed — skipping end-to-end index test")
	}
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "layouts/base.html", `<html><body>{% block content %}{{ content | safe }}{% endblock %}</body></html>`)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\nlayout: base.html\n---\nHello world, this is searchable content.\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	n, err := runSearchIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Error("expected pagefind to write at least one file")
	}
	if _, err := os.Stat(filepath.Join(root, "static", "pagefind", "pagefind.js")); err != nil {
		t.Errorf("static/pagefind/pagefind.js not written: %v", err)
	}

	if err := os.RemoveAll(filepath.Join(root, "site")); err != nil {
		t.Fatal(err)
	}
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "site", "pagefind", "pagefind.js")); err != nil {
		t.Errorf("site/pagefind/pagefind.js not present after a plain build: %v", err)
	}
}
