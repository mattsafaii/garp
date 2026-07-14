package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	writeOut("style.css", "body { color: red }")
	if with404 {
		writeOut("404.html", "<h1>lost</h1>")
	}
	return dir
}

func get(t *testing.T, dir, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	siteHandler(dir, newReloader()).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
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

// TestDevInjectsReloadScript: every HTML response dev serves — clean
// URLs, the index, the 404 page — carries the live-reload script. The
// script exists only in dev responses; site/ on disk stays untouched.
func TestDevInjectsReloadScript(t *testing.T) {
	dir := devFixture(t, true)
	for _, p := range []string{"/", "/about", "/about.html", "/nope"} {
		if body := get(t, dir, p).Body.String(); !strings.Contains(body, "/_garp/reload") {
			t.Errorf("GET %s missing live-reload script:\n%s", p, body)
		}
	}
	if b, err := os.ReadFile(filepath.Join(dir, "about.html")); err != nil || strings.Contains(string(b), "_garp") {
		t.Errorf("about.html on disk should be untouched, got: %s (err %v)", b, err)
	}
}

// TestDevDoesNotInjectNonHTML: assets stream through untouched.
func TestDevDoesNotInjectNonHTML(t *testing.T) {
	dir := devFixture(t, true)
	body := get(t, dir, "/style.css").Body.String()
	if strings.Contains(body, "_garp") {
		t.Errorf("style.css should not carry the reload script:\n%s", body)
	}
}

// TestDevReloadSSE: a broadcast reaches a connected /_garp/reload
// subscriber as an SSE event.
func TestDevReloadSSE(t *testing.T) {
	dir := devFixture(t, false)
	rl := newReloader()
	srv := httptest.NewServer(siteHandler(dir, rl))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/_garp/reload")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	// wait for the subscription to register, then fire
	deadline := time.Now().Add(2 * time.Second)
	for {
		rl.mu.Lock()
		n := len(rl.subs)
		rl.mu.Unlock()
		if n > 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("subscriber never registered")
		}
		time.Sleep(10 * time.Millisecond)
	}
	rl.broadcast()

	got := make(chan string, 1)
	go func() {
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			if line := sc.Text(); line != "" {
				got <- line
				return
			}
		}
	}()
	select {
	case line := <-got:
		if line != "data: reload" {
			t.Errorf("SSE line = %q, want %q", line, "data: reload")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no SSE event within 2s")
	}
}
