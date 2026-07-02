package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyStatic(t *testing.T) {
	src, out := t.TempDir(), t.TempDir()
	writeFile(t, src, "style.css", "body { margin: 0 }")
	writeFile(t, src, "img/logo.png", "\x89PNG fake bytes")
	writeFile(t, src, "_headers", "/*\n  X-Frame-Options: DENY\n")

	n, collisions, err := copyStatic(src, out, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("copied %d files, want 3", n)
	}
	if len(collisions) != 0 {
		t.Errorf("collisions = %v, want none", collisions)
	}
	for _, rel := range []string{"style.css", "img/logo.png", "_headers"} {
		want, err := os.ReadFile(filepath.Join(src, rel))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(out, rel))
		if err != nil {
			t.Fatalf("%s not copied: %v", rel, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s not byte-identical", rel)
		}
	}
}

func TestCopyStaticMissingDir(t *testing.T) {
	n, _, err := copyStatic(filepath.Join(t.TempDir(), "nope"), t.TempDir(), nil)
	if err != nil || n != 0 {
		t.Errorf("n = %d, err = %v", n, err)
	}
}

// A static file landing on a rendered page's output path is reported —
// static wins silently otherwise, which is a confusing debugging session.
func TestCopyStaticReportsPageCollisions(t *testing.T) {
	src, out := t.TempDir(), t.TempDir()
	writeFile(t, src, "index.html", "static wins")
	writeFile(t, src, "style.css", "body{}")

	n, collisions, err := copyStatic(src, out, map[string]bool{"index.html": true})
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("copied %d files, want 2", n)
	}
	if len(collisions) != 1 || collisions[0] != "index.html" {
		t.Errorf("collisions = %v, want [index.html]", collisions)
	}
}
