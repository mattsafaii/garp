package main

import "strings"

// outputPath maps a page to its file under site/: a 1:1 flat-.html mirror of
// content/ (about.md → about.html), unless the cascade sets permalink:.
// Permalinks are written URL-style — "/contact-us" → contact-us.html,
// "/docs/" → docs/index.html, "/feed.html" stays as-is; anything without
// a .html suffix gets one appended.
func outputPath(page *Page, data map[string]any) string {
	p, ok := data["permalink"].(string)
	if !ok || p == "" {
		return strings.TrimSuffix(page.Source, ".md") + ".html"
	}
	p = strings.TrimPrefix(p, "/")
	switch {
	case p == "":
		return "index.html"
	case strings.HasSuffix(p, "/"):
		return p + "index.html"
	case strings.HasSuffix(p, ".html"):
		return p
	default:
		return p + ".html"
	}
}

// pageURL is the clean URL a host like Cloudflare Pages serves the output
// file at: about.html → /about, index.html → /, blog/index.html → /blog.
func pageURL(outPath string) string {
	url := "/" + strings.TrimSuffix(outPath, ".html")
	if url == "/index" {
		return "/"
	}
	return strings.TrimSuffix(url, "/index")
}
