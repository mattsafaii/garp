package main

import (
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	out, err := renderMarkdown("# Hi\n\nSome **bold** text.\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `<h1 id="hi">Hi</h1>`) || !strings.Contains(out, "<strong>bold</strong>") {
		t.Errorf("out = %q", out)
	}
}

func TestRenderMarkdownGFMTable(t *testing.T) {
	out, err := renderMarkdown("| a | b |\n|---|---|\n| 1 | 2 |\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<table>") {
		t.Errorf("GFM table not rendered: %q", out)
	}
}

func TestRenderMarkdownRawHTML(t *testing.T) {
	out, err := renderMarkdown("<div class=\"hero\">raw</div>\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `<div class="hero">`) {
		t.Errorf("raw HTML should pass through: %q", out)
	}
}
