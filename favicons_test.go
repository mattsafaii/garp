package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// squareSource writes a size x size solid-color PNG to dir/name and returns
// its path — a stand-in for the "one square source image" the command reads.
func squareSource(t *testing.T, dir, name string, size int) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return path
}

func faviconsFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	return root
}

func TestGenerateFaviconsWritesFullSet(t *testing.T) {
	root := faviconsFixture(t)
	source := squareSource(t, root, "logo.png", 1024)

	written, err := generateFavicons(root, source)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 5 {
		t.Errorf("wrote %d files, want 5: %v", len(written), written)
	}

	for _, name := range []string{"favicon.ico", "icon-192.png", "icon-512.png", "apple-touch-icon.png", "site.webmanifest"} {
		if _, err := os.Stat(filepath.Join(root, "static", name)); err != nil {
			t.Errorf("static/%s not written: %v", name, err)
		}
	}

	icon192, err := os.Open(filepath.Join(root, "static", "icon-192.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer icon192.Close()
	cfg, err := png.DecodeConfig(icon192)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 192 || cfg.Height != 192 {
		t.Errorf("icon-192.png is %dx%d, want 192x192", cfg.Width, cfg.Height)
	}
}

func TestGenerateFaviconsByteStable(t *testing.T) {
	root1 := faviconsFixture(t)
	source1 := squareSource(t, root1, "logo.png", 1024)
	if _, err := generateFavicons(root1, source1); err != nil {
		t.Fatal(err)
	}

	root2 := faviconsFixture(t)
	source2 := squareSource(t, root2, "logo.png", 1024)
	if _, err := generateFavicons(root2, source2); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"favicon.ico", "icon-192.png", "icon-512.png", "apple-touch-icon.png", "site.webmanifest"} {
		a, err := os.ReadFile(filepath.Join(root1, "static", name))
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(root2, "static", name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, b) {
			t.Errorf("%s not byte-stable across runs", name)
		}
	}
}

func TestGenerateFaviconsRejectsNonSquare(t *testing.T) {
	root := faviconsFixture(t)
	img := image.NewRGBA(image.Rect(0, 0, 1024, 512))
	path := filepath.Join(root, "logo.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	png.Encode(f, img)
	f.Close()

	if _, err := generateFavicons(root, path); err == nil {
		t.Error("expected an error for a non-square source image")
	}
}

func TestGenerateFaviconsRejectsSmallSource(t *testing.T) {
	root := faviconsFixture(t)
	source := squareSource(t, root, "logo.png", 512)
	if _, err := generateFavicons(root, source); err == nil {
		t.Error("expected an error for a source image under 1024px")
	}
}

func TestGenerateFaviconsRequiresConfig(t *testing.T) {
	root := t.TempDir()
	source := squareSource(t, root, "logo.png", 1024)
	if _, err := generateFavicons(root, source); err == nil {
		t.Error("expected an error without config.yaml")
	}
}

func TestGenerateFaviconsManifestUsesSiteName(t *testing.T) {
	root := faviconsFixture(t)
	source := squareSource(t, root, "logo.png", 1024)
	if _, err := generateFavicons(root, source); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "static", "site.webmanifest"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatalf("site.webmanifest is not valid JSON: %v", err)
	}
	if manifest["name"] != "Fixture" {
		t.Errorf("manifest name = %v, want Fixture", manifest["name"])
	}
	icons, _ := manifest["icons"].([]any)
	if len(icons) != 2 {
		t.Errorf("manifest icons = %v, want 2 entries", manifest["icons"])
	}
}

// TestFaviconsPartialRenders exercises the real scaffold favicons.html
// component included from head.html, guarding the acceptance criterion that
// the built site serves correct head links after garp favicons runs.
func TestFaviconsPartialRenders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "components/favicons.html", readScaffold(t, "scaffold/components/favicons.html"))
	writeFile(t, root, "layouts/base.html", `<head>{% include "favicons.html" %}</head>{% block content %}{{ content | safe }}{% endblock %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	html := read(t, root, "index.html")
	for _, want := range []string{
		`<link rel="icon" href="/favicon.ico" sizes="32x32">`,
		`<link rel="icon" type="image/png" href="/icon-192.png" sizes="192x192">`,
		`<link rel="icon" type="image/png" href="/icon-512.png" sizes="512x512">`,
		`<link rel="apple-touch-icon" href="/apple-touch-icon.png">`,
		`<link rel="manifest" href="/site.webmanifest">`,
	} {
		if !bytes.Contains([]byte(html), []byte(want)) {
			t.Errorf("favicons.html missing %q:\n%s", want, html)
		}
	}
}
