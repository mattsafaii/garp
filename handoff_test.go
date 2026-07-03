package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func handoffFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "content/index.md", "home\n")
	writeFile(t, root, "content/about.md", "about\n")
	writeFile(t, root, "data/nav.yaml", "- Home\n")
	return root
}

func TestRunHandoffRequiresConfig(t *testing.T) {
	if _, _, err := runHandoff(t.TempDir()); err == nil {
		t.Error("expected an error without config.yaml")
	}
}

func TestRunHandoffWritesBinAndDoc(t *testing.T) {
	root := handoffFixture(t)
	written, _, err := runHandoff(root)
	if err != nil {
		t.Fatal(err)
	}

	currentPlatform := "garp-" + runtime.GOOS + "-" + runtime.GOARCH
	wantFiles := []string{filepath.Join("bin", currentPlatform), filepath.Join("bin", "garp"), "HANDOFF.md"}
	for _, w := range wantFiles {
		found := false
		for _, got := range written {
			if got == w {
				found = true
			}
		}
		if !found {
			t.Errorf("runHandoff did not report writing %s (wrote %v)", w, written)
		}
		if _, err := os.Stat(filepath.Join(root, w)); err != nil {
			t.Errorf("%s not on disk: %v", w, err)
		}
	}

	fi, err := os.Stat(filepath.Join(root, "bin", "garp"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error("bin/garp shim is not executable")
	}

	shim, err := os.ReadFile(filepath.Join(root, "bin", "garp"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(shim), "#!/bin/sh\n") {
		t.Errorf("bin/garp shim missing shebang:\n%s", shim)
	}
	if !strings.Contains(string(shim), "uname -s") || !strings.Contains(string(shim), "uname -m") {
		t.Errorf("bin/garp shim should dispatch on uname:\n%s", shim)
	}
}

func TestRunHandoffWarnsWithoutCIBinary(t *testing.T) {
	root := handoffFixture(t)
	_, warning, err := runHandoff(root)
	if err != nil {
		t.Fatal(err)
	}
	// under `go test`, the running binary's directory won't contain a
	// garp-linux-amd64 sibling (unless the test happens to run ON linux/amd64
	// itself, in which case there's nothing to warn about)
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		t.Skip("running on the CI platform itself — no warning expected")
	}
	if warning == "" {
		t.Error("expected a warning when garp-linux-amd64 isn't present")
	}
	if !strings.Contains(warning, ciBinary) {
		t.Errorf("warning should mention %s: %s", ciBinary, warning)
	}
}

func TestRunHandoffDocReflectsProjectShape(t *testing.T) {
	root := handoffFixture(t)
	if _, err := stampBlog(root); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "static/favicon.ico", "fake")
	writeFile(t, root, "static/og/index.png", "fake")

	_, _, err := runHandoff(root)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := os.ReadFile(filepath.Join(root, "HANDOFF.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(doc)
	for _, want := range []string{
		"# Fixture — Handoff",
		"https://fixture.test",
		"**Content sections:** blog", // top-level content/ subdirectories only, not top-level files
		"nav",                        // data/nav.yaml
		"**Blog:** yes",
		"**Atom feed:** yes",
		"**Favicons:** yes",
		"**OG images:** yes",
		"bin/garp favicons <source>",
		"bin/garp og",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("HANDOFF.md missing %q:\n%s", want, s)
		}
	}
}

func TestRunHandoffIsRerunnable(t *testing.T) {
	root := handoffFixture(t)
	if _, _, err := runHandoff(root); err != nil {
		t.Fatal(err)
	}
	if _, _, err := runHandoff(root); err != nil {
		t.Errorf("second handoff run should succeed, not refuse: %v", err)
	}
}
