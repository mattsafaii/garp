package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildFixture is a small site exercising the whole pipeline: cascade,
// layouts with extends/include, collections, permalink, static.
func buildFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\nphone: \"555\"\n")
	writeFile(t, root, "data/nav.yaml", "- Home\n- Blog\n")
	writeFile(t, root, "layouts/base.html", `<title>{{ title }} | {{ site.site_name }}</title>
<nav>{% for item in nav %}<a>{{ item }}</a>{% endfor %}</nav>
{% block content %}{{ content | safe }}{% endblock %}
{% include "footer.html" %}`)
	writeFile(t, root, "layouts/post.html", `{% extends "base.html" %}
{% block content %}<article>{{ content | safe }}</article>{% endblock %}`)
	writeFile(t, root, "components/footer.html", `<footer>{{ site.phone }}</footer>`)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\nlayout: base.html\n---\n# Welcome\n\nCall {{ site.phone }}.\n")
	writeFile(t, root, "content/about.md", "---\ntitle: About\nlayout: base.html\npermalink: /about-us\n---\nabout body\n")
	writeFile(t, root, "content/blog/_data.yaml", "layout: post.html\n")
	writeFile(t, root, "content/blog/index.md", `---
title: Blog
layout: base.html
---
{% for post in collections.blog %}* [{{ post.title }}]({{ post.url }})
{% endfor %}`)
	writeFile(t, root, "content/blog/first.md", "---\ntitle: First Post\ndate: 2026-01-05\n---\nfirst\n")
	writeFile(t, root, "content/blog/second.md", "---\ntitle: Second Post\ndate: 2026-02-10\n---\nsecond\n")
	writeFile(t, root, "static/css/style.css", "body{}")
	return root
}

func read(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, "site", rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBuildSite(t *testing.T) {
	root := buildFixture(t)
	n, err := buildSite(root)
	if err != nil {
		t.Fatal(err)
	}
	if n != 6 {
		t.Errorf("built %d files, want 6 (5 pages + 1 static)", n)
	}

	index := read(t, root, "index.html")
	for _, want := range []string{"<title>Home | Fixture</title>", "<h1>Welcome</h1>", "Call 555.", "<a>Home</a>", "<footer>555</footer>"} {
		if !strings.Contains(index, want) {
			t.Errorf("index.html missing %q:\n%s", want, index)
		}
	}

	// permalink override
	if _, err := os.Stat(filepath.Join(root, "site", "about-us.html")); err != nil {
		t.Error("permalink output missing:", err)
	}

	// directory data applies the post layout; extends chains to base
	first := read(t, root, "blog/first.html")
	if !strings.Contains(first, "<article><p>first</p>") {
		t.Errorf("post layout not applied via dir data:\n%s", first)
	}
	if !strings.Contains(first, "<title>First Post | Fixture</title>") {
		t.Errorf("extends chain broken:\n%s", first)
	}

	// blog index loops collections.blog, newest first, excluding itself
	blogIndex := read(t, root, "blog/index.html")
	if strings.Contains(blogIndex, "Blog</a>") && strings.Contains(blogIndex, "/blog\"") {
		t.Errorf("blog index should not list itself:\n%s", blogIndex)
	}
	si := strings.Index(blogIndex, "Second Post")
	fi := strings.Index(blogIndex, "First Post")
	if si == -1 || fi == -1 || si > fi {
		t.Errorf("collection order wrong (want Second before First):\n%s", blogIndex)
	}
	if !strings.Contains(blogIndex, "/blog/first") {
		t.Errorf("collection urls missing:\n%s", blogIndex)
	}

	// static passthrough
	if read(t, root, "css/style.css") != "body{}" {
		t.Error("static file not copied verbatim")
	}
}

func TestBuildSiteRemovesStaleOutput(t *testing.T) {
	root := buildFixture(t)
	writeFile(t, root, "site/stale.html", "old")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "site", "stale.html")); err == nil {
		t.Error("stale output should be removed by rebuild")
	}
}

func TestBuildSiteNoConfig(t *testing.T) {
	_, err := buildSite(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "config.yaml") {
		t.Errorf("err = %v", err)
	}
}
