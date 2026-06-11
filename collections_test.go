package main

import (
	"testing"
	"time"
)

func ref(source string, data map[string]any) *pageRef {
	page := &Page{Source: source, Front: data}
	out := outputPath(page, data)
	return &pageRef{page: page, data: data, out: out, url: pageURL(out)}
}

func entries(t *testing.T, collections map[string]any, name string) []map[string]any {
	t.Helper()
	v, ok := collections[name].([]map[string]any)
	if !ok {
		t.Fatalf("collections[%q] = %#v", name, collections[name])
	}
	return v
}

func TestCollections(t *testing.T) {
	d1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	cols := buildCollections([]*pageRef{
		ref("index.md", map[string]any{}),
		ref("blog/index.md", map[string]any{}),
		ref("blog/old.md", map[string]any{"date": d1, "title": "Old"}),
		ref("blog/new.md", map[string]any{"date": d2, "title": "New"}),
		ref("services.md", map[string]any{"tags": "featured"}),
	})

	blog := entries(t, cols, "blog")
	if len(blog) != 2 {
		t.Fatalf("blog has %d entries, want 2 (index.md excluded): %v", len(blog), blog)
	}
	if blog[0]["title"] != "New" || blog[1]["title"] != "Old" {
		t.Errorf("blog order = %v, %v — want newest first", blog[0]["title"], blog[1]["title"])
	}
	if blog[0]["url"] != "/blog/new" {
		t.Errorf("url = %v", blog[0]["url"])
	}

	featured := entries(t, cols, "featured")
	if len(featured) != 1 || featured[0]["url"] != "/services" {
		t.Errorf("featured = %v", featured)
	}

	if len(entries(t, cols, "all")) != 5 {
		t.Errorf("all = %d entries", len(entries(t, cols, "all")))
	}
}

func TestCollectionsTagListAndDedupe(t *testing.T) {
	cols := buildCollections([]*pageRef{
		ref("blog/post.md", map[string]any{"tags": []any{"blog", "go"}}),
	})
	if got := len(entries(t, cols, "blog")); got != 1 {
		t.Errorf("blog tag + dir grouping should dedupe, got %d entries", got)
	}
	if got := len(entries(t, cols, "go")); got != 1 {
		t.Errorf("go = %d entries", got)
	}
}
