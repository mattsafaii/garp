package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// devFixture writes a minimal output dir for siteHandler tests.
func devFixture(t *testing.T, with404 bool) string {
	t.Helper()
	dir := t.TempDir()
	writeOut := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeOut("index.html", "<h1>home</h1>")
	writeOut("about.html", "<h1>about</h1>")
	if with404 {
		writeOut("404.html", "<h1>lost</h1>")
	}
	return dir
}

func get(t *testing.T, dir, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	siteHandler(dir).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
	return rec
}

// TestDevCleanURLFallback guards the host-mirroring contract: /about
// serves about.html.
func TestDevCleanURLFallback(t *testing.T) {
	dir := devFixture(t, true)
	rec := get(t, dir, "/about")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "<h1>about</h1>") {
		t.Errorf("GET /about = %d %q, want 200 with about.html", rec.Code, rec.Body.String())
	}
}

// TestDevServes404Page mirrors Cloudflare Pages: a missing path gets the
// project's 404.html with a 404 status.
func TestDevServes404Page(t *testing.T) {
	dir := devFixture(t, true)
	rec := get(t, dir, "/nope")
	if rec.Code != 404 {
		t.Errorf("GET /nope status = %d, want 404", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<h1>lost</h1>") {
		t.Errorf("GET /nope body = %q, want 404.html content", rec.Body.String())
	}
}

// TestDevMissingWithout404Page keeps the fallback sane when a project has
// no 404.html: still a 404, just FileServer's plain one.
func TestDevMissingWithout404Page(t *testing.T) {
	dir := devFixture(t, false)
	rec := get(t, dir, "/nope")
	if rec.Code != 404 {
		t.Errorf("GET /nope status = %d, want 404", rec.Code)
	}
}
