package main

import "testing"

func TestOutputPath(t *testing.T) {
	cases := []struct {
		source    string
		permalink string
		want      string
	}{
		{"index.md", "", "index.html"},
		{"about.md", "", "about.html"},
		{"blog/post.md", "", "blog/post.html"},
		{"blog/index.md", "", "blog/index.html"},
		{"about.md", "/about-us", "about-us.html"},
		{"about.md", "/docs/", "docs/index.html"},
		{"about.md", "/legal/terms.html", "legal/terms.html"},
		{"about.md", "/feed.xml", "feed.xml"},
		{"about.md", "/humans.txt", "humans.txt"},
		{"about.md", "/", "index.html"},
	}
	for _, c := range cases {
		data := map[string]any{}
		if c.permalink != "" {
			data["permalink"] = c.permalink
		}
		got := outputPath(&Page{Source: c.source}, data)
		if got != c.want {
			t.Errorf("outputPath(%q, permalink=%q) = %q, want %q", c.source, c.permalink, got, c.want)
		}
	}
}

func TestPageURL(t *testing.T) {
	cases := map[string]string{
		"index.html":      "/",
		"about.html":      "/about",
		"blog/post.html":  "/blog/post",
		"blog/index.html": "/blog",
	}
	for out, want := range cases {
		if got := pageURL(out); got != want {
			t.Errorf("pageURL(%q) = %q, want %q", out, got, want)
		}
	}
}
