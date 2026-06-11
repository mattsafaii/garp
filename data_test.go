package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadGlobalData(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "nav.yaml", "- Home\n- About\n")
	writeFile(t, dir, "company.json", `{"name": "Zonebrite"}`)
	writeFile(t, dir, "readme.txt", "ignored")

	global, err := loadGlobalData(dir)
	if err != nil {
		t.Fatal(err)
	}
	nav, ok := global["nav"].([]any)
	if !ok || len(nav) != 2 {
		t.Errorf("nav = %v", global["nav"])
	}
	company, ok := global["company"].(map[string]any)
	if !ok || company["name"] != "Zonebrite" {
		t.Errorf("company = %v", global["company"])
	}
	if _, ok := global["readme"]; ok {
		t.Error("non-data file should be ignored")
	}
}

func TestLoadGlobalDataMissingDir(t *testing.T) {
	global, err := loadGlobalData(filepath.Join(t.TempDir(), "nope"))
	if err != nil {
		t.Fatal(err)
	}
	if len(global) != 0 {
		t.Errorf("global = %v", global)
	}
}

func TestCascadePageWins(t *testing.T) {
	content := t.TempDir()
	writeFile(t, content, "blog/_data.yaml", "layout: post.html\nauthor: Matt\n")

	global := map[string]any{"author": "Global", "color": "blue"}
	dirData, err := loadDirData(content)
	if err != nil {
		t.Fatal(err)
	}

	page := &Page{Source: "blog/post.md", Front: map[string]any{"author": "Page"}}
	data := pageData(global, dirData, page)

	if data["author"] != "Page" {
		t.Errorf("author = %v, want frontmatter to win", data["author"])
	}
	if data["layout"] != "post.html" {
		t.Errorf("layout = %v, want directory data", data["layout"])
	}
	if data["color"] != "blue" {
		t.Errorf("color = %v, want global", data["color"])
	}

	// A page outside blog/ gets no directory data.
	other := pageData(global, dirData, &Page{Source: "index.md", Front: map[string]any{}})
	if _, ok := other["layout"]; ok {
		t.Error("root page should not inherit blog/_data.yaml")
	}
}
