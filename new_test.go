package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewScaffold runs `garp new` into a temp dir and asserts every embedded
// scaffold file is emitted byte-for-byte, plus the reserved dirs exist. This is
// the golden test for the scaffolder: it guards the embed+copy path so a
// mangled or dropped file can't ship silently.
func TestNewScaffold(t *testing.T) {
	root := filepath.Join(t.TempDir(), "mysite")
	if err := cmdNew([]string{root}); err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{"content", "layouts", "components", "data", "static", "site"} {
		if fi, err := os.Stat(filepath.Join(root, dir)); err != nil || !fi.IsDir() {
			t.Errorf("reserved dir %s missing", dir)
		}
	}

	err := fs.WalkDir(scaffoldFS, "scaffold", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel("scaffold", path)
		if err != nil {
			return err
		}
		want, err := scaffoldFS.ReadFile(path)
		if err != nil {
			return err
		}
		got, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Errorf("%s not emitted: %v", rel, err)
			return nil
		}
		if string(got) != string(want) {
			t.Errorf("%s not byte-identical to embedded scaffold", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestNewRejectsExisting confirms `garp new` refuses a path that already exists.
func TestNewRejectsExisting(t *testing.T) {
	root := t.TempDir() // already exists
	if err := cmdNew([]string{root}); err == nil {
		t.Error("expected error for existing path, got nil")
	}
}

// TestScaffoldContract encodes the CSS↔HTML contract the scaffold depends on:
// the starter stylesheet lays content out in a grid where only <main>'s
// children reach the content column, so base.html MUST wrap content in <main>.
// This guards against syncing the CSS while leaving the layout behind — the
// exact drift that shipped a broken scaffold before this test existed.
func TestScaffoldContract(t *testing.T) {
	css := readScaffold(t, "scaffold/static/style.css")
	layout := readScaffold(t, "scaffold/layouts/base.html")

	if strings.Contains(css, "grid-column: content") && !strings.Contains(layout, "<main>") {
		t.Error("style.css places children in a content grid but base.html has no <main> wrapper")
	}
	if !strings.Contains(css, "@layer reset, tokens, base, composition, block, utility, exception;") {
		t.Error("style.css missing the canonical @layer declaration — looks empty or mangled")
	}
}

func readScaffold(t *testing.T, path string) string {
	t.Helper()
	b, err := scaffoldFS.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
