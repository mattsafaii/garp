package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	front, body, err := parseFrontmatter([]byte("---\ntitle: Hi\ntags: [blog]\n---\n# Heading\n"))
	if err != nil {
		t.Fatal(err)
	}
	if front["title"] != "Hi" {
		t.Errorf("title = %v", front["title"])
	}
	if body != "# Heading\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatterAbsent(t *testing.T) {
	front, body, err := parseFrontmatter([]byte("# Just markdown\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(front) != 0 {
		t.Errorf("front = %v, want empty", front)
	}
	if body != "# Just markdown\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatterEmptyBody(t *testing.T) {
	front, body, err := parseFrontmatter([]byte("---\ntitle: Hi\n---"))
	if err != nil {
		t.Fatal(err)
	}
	if front["title"] != "Hi" || body != "" {
		t.Errorf("front = %v, body = %q", front, body)
	}
}

func TestParseFrontmatterEmpty(t *testing.T) {
	for src, wantBody := range map[string]string{
		"---\n---\nbody\n": "body\n",
		"---\n---\n":       "",
		"---\n---":         "",
	} {
		front, body, err := parseFrontmatter([]byte(src))
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		if len(front) != 0 || body != wantBody {
			t.Errorf("%q: front = %v, body = %q, want empty front, body %q", src, front, body, wantBody)
		}
	}
}

// A leading --- that never closes isn't frontmatter — e.g. a file opening
// with a thematic break. The whole source is body.
func TestParseFrontmatterUnclosedIsBody(t *testing.T) {
	for _, src := range []string{"---\n\na horizontal rule opener\n", "---\ntitle: Hi\n"} {
		front, body, err := parseFrontmatter([]byte(src))
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		if len(front) != 0 || body != src {
			t.Errorf("%q: front = %v, body = %q, want full source as body", src, front, body)
		}
	}
}

func TestParseFrontmatterHorizontalRuleInBody(t *testing.T) {
	_, body, err := parseFrontmatter([]byte("---\ntitle: Hi\n---\nabove\n\n---\n\nbelow\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "below") {
		t.Errorf("body lost content after hr: %q", body)
	}
}

func TestDiscoverContent(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("index.md", "---\ntitle: Home\n---\nhi\n")
	write("blog/post.md", "no frontmatter\n")
	write("notes.txt", "not a page")
	write(".hidden.md", "skip me")

	pages, err := discoverContent(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2: %+v", len(pages), pages)
	}
	if pages[0].Source != "blog/post.md" || pages[1].Source != "index.md" {
		t.Errorf("sources = %q, %q", pages[0].Source, pages[1].Source)
	}
	if pages[1].Front["title"] != "Home" {
		t.Errorf("front = %v", pages[1].Front)
	}
}
