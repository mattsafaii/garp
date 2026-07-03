package main

import (
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func blogFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "layouts/base.html", `<title>{{ title }}</title>{% block content %}{{ content | safe }}{% endblock %}`)
	return root
}

func TestStampBlogWritesExpectedFiles(t *testing.T) {
	root := blogFixture(t)
	written, err := stampBlog(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"content/blog/_data.yaml",
		"content/blog/hello-world.md",
		"content/blog/index.md",
		"content/feed.md",
		"layouts/feed.xml",
		"layouts/post.html",
	}
	if len(written) != len(want) {
		t.Fatalf("wrote %v, want %v", written, want)
	}
	for i, w := range want {
		if written[i] != w {
			t.Errorf("written[%d] = %q, want %q", i, written[i], w)
		}
		if _, err := os.Stat(filepath.Join(root, w)); err != nil {
			t.Errorf("%s not written: %v", w, err)
		}
	}
}

func TestStampBlogRefusesToOverwrite(t *testing.T) {
	root := blogFixture(t)
	if _, err := stampBlog(root); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, "content/blog/hello-world.md"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := stampBlog(root); err == nil {
		t.Error("expected an error running garp blog twice")
	}

	after, err := os.ReadFile(filepath.Join(root, "content/blog/hello-world.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("second run should not have touched existing files")
	}
}

// TestStampBlogGolden guards the embed+copy path itself: every file stamped
// into a project must be byte-identical to the embedded scaffold-blog/
// source, the same guarantee TestNewScaffold gives scaffold/.
func TestStampBlogGolden(t *testing.T) {
	root := blogFixture(t)
	if _, err := stampBlog(root); err != nil {
		t.Fatal(err)
	}
	err := fs.WalkDir(blogFS, "scaffold-blog", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel("scaffold-blog", path)
		if err != nil {
			return err
		}
		want, err := blogFS.ReadFile(path)
		if err != nil {
			return err
		}
		got, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Errorf("%s not emitted: %v", rel, err)
			return nil
		}
		if string(got) != string(want) {
			t.Errorf("%s not byte-identical to embedded scaffold-blog", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStampBlogRefusesOnExistingProjectFile(t *testing.T) {
	root := blogFixture(t)
	writeFile(t, root, "layouts/post.html", "custom post layout\n")
	if _, err := stampBlog(root); err == nil {
		t.Error("expected an error when a target file already exists")
	}
	// nothing else should have been written either — all or nothing
	if _, err := os.Stat(filepath.Join(root, "content/blog/_data.yaml")); err == nil {
		t.Error("stampBlog should write nothing when any target file conflicts")
	}
}

// TestBlogEndToEnd runs the stamped scaffold through a real build: the
// listing page lists the sample post, the post renders via the post layout,
// and the Atom feed is well-formed with the required elements.
func TestBlogEndToEnd(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "layouts/base.html", `<title>{% if title %}{{ title }} — {% endif %}Fixture</title>{% block content %}{{ content | safe }}{% endblock %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	if _, err := stampBlog(root); err != nil {
		t.Fatal(err)
	}
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	index := read(t, root, "blog/index.html")
	if !strings.Contains(index, `href="/blog/hello-world"`) || !strings.Contains(index, "Hello, World") {
		t.Errorf("blog index missing the sample post:\n%s", index)
	}

	post := read(t, root, "blog/hello-world.html")
	if !strings.Contains(post, "<h1>Hello, World</h1>") || !strings.Contains(post, "This is your first post") {
		t.Errorf("post did not render via the post layout:\n%s", post)
	}

	feed := read(t, root, "feed.xml")
	var parsed struct {
		XMLName xml.Name `xml:"feed"`
		Title   string   `xml:"title"`
		ID      string   `xml:"id"`
		Updated string   `xml:"updated"`
		Entries []struct {
			Title   string `xml:"title"`
			ID      string `xml:"id"`
			Updated string `xml:"updated"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal([]byte(feed), &parsed); err != nil {
		t.Fatalf("feed.xml did not parse as XML: %v\n%s", err, feed)
	}
	if parsed.Title != "Fixture" {
		t.Errorf("feed title = %q, want Fixture", parsed.Title)
	}
	if parsed.ID == "" || parsed.Updated == "" {
		t.Errorf("feed missing required id/updated: %+v", parsed)
	}
	if len(parsed.Entries) != 1 {
		t.Fatalf("feed has %d entries, want 1", len(parsed.Entries))
	}
	e := parsed.Entries[0]
	if e.Title != "Hello, World" || e.ID != "https://fixture.test/blog/hello-world" || e.Updated == "" {
		t.Errorf("feed entry wrong: %+v", e)
	}
}
