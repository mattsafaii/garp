package main

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

// ogFixture writes a minimal project with a real committed font (Go's own
// vendored regular face, reused here as a stand-in for a client-provided
// .ttf) and og configured, so generateOGImages exercises real font parsing.
func ogFixture(t *testing.T, ogConfig string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "fonts/test.ttf", string(goregular.TTF))
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n"+ogConfig)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\n---\nhome\n")
	writeFile(t, root, "content/about.md", "---\ntitle: A Fairly Long About Page Title That Should Wrap\n---\nabout\n")
	writeFile(t, root, "content/404.md", "---\ntitle: Not Found\n---\nnot found\n")
	writeFile(t, root, "content/feed.md", "---\npermalink: /feed.xml\n---\nfeed\n")
	return root
}

func TestGenerateOGImagesRequiresFont(t *testing.T) {
	root := ogFixture(t, "")
	if _, err := generateOGImages(root); err == nil {
		t.Error("expected an error when og.font is unset")
	}
}

func TestGenerateOGImagesWritesPNGs(t *testing.T) {
	root := ogFixture(t, "og:\n  font: fonts/test.ttf\n")
	written, err := generateOGImages(root)
	if err != nil {
		t.Fatal(err)
	}
	// index.md and about.md only — 404 and the feed permalink are skipped
	if len(written) != 2 {
		t.Errorf("wrote %d images, want 2: %v", len(written), written)
	}

	for _, rel := range []string{filepath.Join("static", "og", "index.png"), filepath.Join("static", "og", "about.png")} {
		f, err := os.Open(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("%s not written: %v", rel, err)
		}
		cfg, err := png.DecodeConfig(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Width != ogWidth || cfg.Height != ogHeight {
			t.Errorf("%s is %dx%d, want %dx%d", rel, cfg.Width, cfg.Height, ogWidth, ogHeight)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "static", "og", "404.png")); err == nil {
		t.Error("404 page should not get an OG image")
	}
	if _, err := os.Stat(filepath.Join(root, "static", "og", "feed.xml.png")); err == nil {
		t.Error("non-HTML permalink should not get an OG image")
	}
}

func TestGenerateOGImagesByteStable(t *testing.T) {
	root1 := ogFixture(t, "og:\n  font: fonts/test.ttf\n")
	if _, err := generateOGImages(root1); err != nil {
		t.Fatal(err)
	}
	root2 := ogFixture(t, "og:\n  font: fonts/test.ttf\n")
	if _, err := generateOGImages(root2); err != nil {
		t.Fatal(err)
	}

	a, err := os.ReadFile(filepath.Join(root1, "static", "og", "about.png"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root2, "static", "og", "about.png"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Error("OG image not byte-stable across runs")
	}
}

func TestGenerateOGImagesCustomBackground(t *testing.T) {
	root := ogFixture(t, "og:\n  font: fonts/test.ttf\n  background: \"#ff0000\"\n")
	if _, err := generateOGImages(root); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(filepath.Join(root, "static", "og", "index.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := img.At(2, 2).RGBA()
	if r>>8 != 0xff || g>>8 != 0x00 || b>>8 != 0x00 {
		t.Errorf("corner pixel = %d,%d,%d, want red background", r>>8, g>>8, b>>8)
	}
}

func TestGenerateOGImagesBadBackground(t *testing.T) {
	root := ogFixture(t, "og:\n  font: fonts/test.ttf\n  background: not-a-color\n")
	if _, err := generateOGImages(root); err == nil {
		t.Error("expected an error for an invalid og.background")
	}
}

// TestOGFallbackMetaTag exercises the real scaffold head.html against the
// acceptance criteria: with site.og set, a page with no image: gets the
// conventional-path fallback, and an explicit image: still overrides it.
func TestOGFallbackMetaTag(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\nog:\n  font: fonts/test.ttf\n")
	writeFile(t, root, "components/head.html", readScaffold(t, "scaffold/components/head.html"))
	writeFile(t, root, "components/jsonld.html", readScaffold(t, "scaffold/components/jsonld.html"))
	writeFile(t, root, "components/favicons.html", readScaffold(t, "scaffold/components/favicons.html"))
	writeFile(t, root, "components/prefetch.html", readScaffold(t, "scaffold/components/prefetch.html"))
	writeFile(t, root, "layouts/base.html", `<head>{% include "head.html" %}</head>{% block content %}{{ content | safe }}{% endblock %}`)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\nlayout: base.html\n---\nbody\n")
	writeFile(t, root, "content/about.md", "---\ntitle: About\nlayout: base.html\nimage: /custom-og.jpg\n---\nbody\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	index := read(t, root, "index.html")
	if !bytesContains(index, `<meta property="og:image" content="https://fixture.test/og/index.png">`) {
		t.Errorf("index.html missing conventional-path og:image fallback:\n%s", index)
	}
	if !bytesContains(index, `<meta name="twitter:card" content="summary_large_image">`) {
		t.Errorf("index.html twitter:card should be summary_large_image when site.og is set:\n%s", index)
	}

	about := read(t, root, "about.html")
	if !bytesContains(about, `<meta property="og:image" content="https://fixture.test/custom-og.jpg">`) {
		t.Errorf("about.html should keep its explicit image: override:\n%s", about)
	}
	if bytesContains(about, "/og/about.png") {
		t.Errorf("about.html should not fall back to the conventional path when image: is set:\n%s", about)
	}
}

func bytesContains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
