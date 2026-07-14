package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// devFixture writes a minimal project root with a pre-built site/ dir
// for siteHandler tests.
func devFixture(t *testing.T, with404 bool) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\n")
	writeFile(t, root, "site/index.html", "<h1>home</h1>")
	writeFile(t, root, "site/about.html", "<h1>about</h1>")
	writeFile(t, root, "site/style.css", "body { color: red }")
	if with404 {
		writeFile(t, root, "site/404.html", "<h1>lost</h1>")
	}
	return root
}

func get(t *testing.T, root, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	siteHandler(root, filepath.Join(root, "site"), newReloader()).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
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

// TestDevInjectsToolbarScript: every HTML response dev serves — clean
// URLs, the index, the 404 page — carries the toolbar script tag. The
// tag exists only in dev responses; site/ on disk stays untouched.
func TestDevInjectsToolbarScript(t *testing.T) {
	dir := devFixture(t, true)
	for _, p := range []string{"/", "/about", "/about.html", "/nope"} {
		if body := get(t, dir, p).Body.String(); !strings.Contains(body, `<script src="/_garp/toolbar.js" defer></script>`) {
			t.Errorf("GET %s missing toolbar script tag:\n%s", p, body)
		}
	}
	if b, err := os.ReadFile(filepath.Join(dir, "site", "about.html")); err != nil || strings.Contains(string(b), "_garp") {
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
	srv := httptest.NewServer(siteHandler(dir, filepath.Join(dir, "site"), rl))
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

// inspectorFixture writes a project source tree exercising the whole
// cascade for /_garp/page: a global key, a dir-data key shadowing a
// global, a frontmatter key shadowing dir data, tags, and a layout that
// extends another.
func inspectorFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "data/nav.yaml", "items:\n  - Home\n")
	writeFile(t, root, "data/author.yaml", "name: Global Author\n")
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}`)
	writeFile(t, root, "layouts/post.html", `{% extends "base.html" %}`)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\n---\nhome\n")
	writeFile(t, root, "content/blog/_data.yaml", "layout: post.html\nauthor: Dir Author\nsection: news\n")
	writeFile(t, root, "content/blog/index.md", "---\ntitle: Blog\n---\nlisting\n")
	writeFile(t, root, "content/blog/post.md", "---\ntitle: Post\ndate: 2026-05-22\nsection: overridden\ntags: [featured]\n---\nbody\n")
	return root
}

// getPage fetches /_garp/page for path and decodes the JSON response.
func getPage(t *testing.T, root, path string) (int, pageInfo) {
	t.Helper()
	rec := get(t, root, "/_garp/page?path="+path)
	var info pageInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &info); err != nil {
		t.Fatalf("GET /_garp/page?path=%s: bad JSON %q: %v", path, rec.Body.String(), err)
	}
	return rec.Code, info
}

// TestDevPageInspector: the metadata endpoint reports a page's source,
// layout chain (following extends), collections, and the resolved
// cascade with each key labeled by the layer that won it.
func TestDevPageInspector(t *testing.T) {
	root := inspectorFixture(t)
	code, info := getPage(t, root, "/blog/post")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if info.URL != "/blog/post" || info.Source != "content/blog/post.md" {
		t.Errorf("url/source = %q %q, want /blog/post content/blog/post.md", info.URL, info.Source)
	}
	if got := strings.Join(info.LayoutChain, ","); got != "post.html,base.html" {
		t.Errorf("layout_chain = %q, want post.html,base.html", got)
	}
	if got := strings.Join(info.Collections, ","); got != "all,blog,featured" {
		t.Errorf("collections = %q, want all,blog,featured", got)
	}

	sources := map[string]string{}
	values := map[string]any{}
	for _, e := range info.Data {
		sources[e.Key] = e.Source
		values[e.Key] = e.Value
	}
	want := map[string]string{
		"nav":     "data/nav.yaml",           // global, unshadowed
		"author":  "content/blog/_data.yaml", // dir data shadows data/author.yaml
		"layout":  "content/blog/_data.yaml", // dir data only
		"section": "frontmatter",             // frontmatter shadows dir data
		"title":   "frontmatter",
		"date":    "frontmatter",
		"tags":    "frontmatter",
		"site":    "config.yaml",
	}
	for k, src := range want {
		if sources[k] != src {
			t.Errorf("data[%s].source = %q, want %q", k, sources[k], src)
		}
	}
	if values["section"] != "overridden" || values["author"] != "Dir Author" || values["date"] != "2026-05-22" {
		t.Errorf("winning values wrong: section=%v author=%v date=%v", values["section"], values["author"], values["date"])
	}
	if site, ok := values["site"].(map[string]any); !ok || site["site_name"] != "Fixture" {
		t.Errorf("site value = %v, want config map with site_name", values["site"])
	}

	// keys sorted for a stable response
	for i := 1; i < len(info.Data); i++ {
		if info.Data[i-1].Key >= info.Data[i].Key {
			t.Errorf("data not sorted by key: %q before %q", info.Data[i-1].Key, info.Data[i].Key)
		}
	}
}

// TestDevPageInspectorUnknownPath: a URL no content page produces gets a
// JSON 404.
func TestDevPageInspectorUnknownPath(t *testing.T) {
	root := inspectorFixture(t)
	rec := get(t, root, "/_garp/page?path=/nope")
	if rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	var e map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil || e["error"] != "no page at /nope" {
		t.Errorf("body = %q, want {\"error\": \"no page at /nope\"} (err %v)", rec.Body.String(), err)
	}
}

// TestDevServesToolbarJS: the embedded toolbar serves as JS at its fixed
// path, whatever its contents.
func TestDevServesToolbarJS(t *testing.T) {
	dir := devFixture(t, false)
	rec := get(t, dir, "/_garp/toolbar.js")
	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/javascript; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/javascript; charset=utf-8", ct)
	}
	if rec.Body.Len() == 0 {
		t.Error("toolbar.js served empty")
	}
}
