package main

import (
	"strings"
	"testing"
)

// setupTemplates builds a layouts/ + components/ pair exercising extends,
// block, and include, and returns a renderer over them.
func setupTemplates(t *testing.T) *renderer {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "layouts/base.html", `<html><head><title>{{ site.site_name }}</title></head>
<body>{% block content %}{{ content | safe }}{% endblock %}{% include "footer.html" %}</body></html>`)
	writeFile(t, dir, "layouts/post.html", `{% extends "base.html" %}
{% block content %}<article>{{ content | safe }}</article>{% endblock %}`)
	writeFile(t, dir, "components/footer.html", `<footer>{{ site.site_name }} footer</footer>`)
	return newRenderer(dir+"/layouts", dir+"/components")
}

func ctx(layout string) map[string]any {
	m := map[string]any{"site": map[string]any{"site_name": "Test Site"}}
	if layout != "" {
		m["layout"] = layout
	}
	return m
}

func TestRenderPageWithLayout(t *testing.T) {
	r := setupTemplates(t)
	out, err := r.renderPage(&Page{Source: "index.md", Body: "# Hello\n"}, ctx("base.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<title>Test Site</title>", `<h1 id="hello">Hello</h1>`, "<footer>Test Site footer</footer>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestRenderPageExtendsChain(t *testing.T) {
	r := setupTemplates(t)
	out, err := r.renderPage(&Page{Source: "post.md", Body: "body text\n"}, ctx("post.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<article><p>body text</p>") {
		t.Errorf("post block not applied:\n%s", out)
	}
	if !strings.Contains(out, "<title>Test Site</title>") {
		t.Errorf("base layout not inherited:\n%s", out)
	}
}

func TestRenderPageNoLayout(t *testing.T) {
	r := setupTemplates(t)
	out, err := r.renderPage(&Page{Source: "raw.md", Body: "plain\n"}, ctx(""))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "<p>plain</p>" {
		t.Errorf("out = %q", out)
	}
}

func TestRenderPageBodyTemplating(t *testing.T) {
	r := setupTemplates(t)
	out, err := r.renderPage(&Page{Source: "index.md", Body: "Call {{ site.site_name }}\n"}, ctx(""))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Call Test Site") {
		t.Errorf("site.* not resolved in body: %q", out)
	}
}

func TestRenderPageMissingLayout(t *testing.T) {
	r := setupTemplates(t)
	_, err := r.renderPage(&Page{Source: "x.md", Body: "hi"}, ctx("nope.html"))
	if err == nil || !strings.Contains(err.Error(), "nope.html") {
		t.Errorf("err = %v, want missing-layout error naming the layout", err)
	}
}
