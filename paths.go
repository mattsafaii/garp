package main

import (
	"path"
	"strings"
)

// outputPath maps a page to its file under site/: a 1:1 flat-.html mirror of
// content/ (about.md → about.html), unless the cascade sets permalink:.
// Permalinks are written URL-style — "/contact-us" → contact-us.html,
// "/docs/" → docs/index.html. A permalink with an extension is kept as-is
// ("/feed.xml", "/legal/terms.html"), so non-HTML outputs are possible;
// only an extensionless one gets .html appended.
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
	case path.Ext(p) != "":
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
