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

	n, err := copyStatic(src, out)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("copied %d files, want 3", n)
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
	n, err := copyStatic(filepath.Join(t.TempDir(), "nope"), t.TempDir())
	if err != nil || n != 0 {
		t.Errorf("n = %d, err = %v", n, err)
	}
}
