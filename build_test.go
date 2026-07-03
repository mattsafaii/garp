package main

import (
	"encoding/json"
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
{% for post in collections.blog %}* [{{ post.title }}]({{ post.url }}) {{ post.date }} / {{ post.date | date:"Jan 2, 2006" }}
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
	if n != 8 {
		t.Errorf("built %d files, want 8 (5 pages + 1 static + sitemap.xml + robots.txt)", n)
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
	// dates render clean by default (ISO) and via the date filter
	if !strings.Contains(blogIndex, "2026-02-10 / Feb 10, 2026") {
		t.Errorf("date formatting wrong (want clean default + filter):\n%s", blogIndex)
	}

	// static passthrough
	if read(t, root, "css/style.css") != "body{}" {
		t.Error("static file not copied verbatim")
	}
}

func TestBuildSiteSynthesizesSEO(t *testing.T) {
	root := buildFixture(t)
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	sitemap := read(t, root, "sitemap.xml")
	if !strings.Contains(sitemap, "<loc>https://fixture.test/blog/first</loc>") {
		t.Errorf("sitemap missing page url:\n%s", sitemap)
	}
	if !strings.Contains(sitemap, "<lastmod>2026-01-05</lastmod>") {
		t.Errorf("sitemap missing lastmod for dated page:\n%s", sitemap)
	}
	start := strings.Index(sitemap, "<loc>https://fixture.test/about-us</loc>")
	if start == -1 {
		t.Fatalf("sitemap missing undated page url:\n%s", sitemap)
	}
	end := start + strings.Index(sitemap[start:], "</url>")
	if strings.Contains(sitemap[start:end], "<lastmod>") {
		t.Errorf("undated page should not have lastmod:\n%s", sitemap[start:end])
	}

	robots := read(t, root, "robots.txt")
	if !strings.Contains(robots, "Sitemap: https://fixture.test/sitemap.xml") {
		t.Errorf("robots.txt missing sitemap line:\n%s", robots)
	}
}

func TestBuildSiteSitemapExcludes404(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "layouts/base.html", "{% block content %}{{ content | safe }}{% endblock %}")
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	writeFile(t, root, "content/404.md", "---\nlayout: base.html\n---\nnot found\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	sitemap := read(t, root, "sitemap.xml")
	if strings.Contains(sitemap, "/404") {
		t.Errorf("sitemap should exclude 404 page:\n%s", sitemap)
	}
}

// A permalink with a non-.html extension is kept as-is — the page renders to
// feed.xml, not feed.xml.html — and stays out of the sitemap (it isn't a page).
func TestBuildSiteNonHTMLPermalink(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "content/index.md", "home\n")
	writeFile(t, root, "content/feed.md", "---\npermalink: /feed.xml\n---\nfeed body\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "site", "feed.xml")); err != nil {
		t.Error("feed.xml not written:", err)
	}
	if _, err := os.Stat(filepath.Join(root, "site", "feed.xml.html")); err == nil {
		t.Error("permalink extension should be preserved, got feed.xml.html")
	}
	if sitemap := read(t, root, "sitemap.xml"); strings.Contains(sitemap, "feed.xml</loc>") {
		t.Errorf("sitemap should exclude non-html outputs:\n%s", sitemap)
	}
}

// Without base_url the sitemap/robots URLs would be relative — invalid per
// both specs — so synthesis is skipped entirely.
func TestBuildSiteSkipsSEOWithoutBaseURL(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\n")
	writeFile(t, root, "content/index.md", "home\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"sitemap.xml", "robots.txt"} {
		if _, err := os.Stat(filepath.Join(root, "site", f)); err == nil {
			t.Errorf("%s should not be synthesized without base_url", f)
		}
	}
}

func TestBuildSiteSEOOverride(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "layouts/base.html", "{% block content %}{{ content | safe }}{% endblock %}")
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	writeFile(t, root, "static/sitemap.xml", "custom sitemap\n")
	writeFile(t, root, "static/robots.txt", "custom robots\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if read(t, root, "sitemap.xml") != "custom sitemap\n" {
		t.Error("author's sitemap.xml should win over synthesis")
	}
	if read(t, root, "robots.txt") != "custom robots\n" {
		t.Error("author's robots.txt should win over synthesis")
	}
}

// analyticsFixture builds a minimal site including the real scaffold
// analytics.html, so these tests exercise the shipped file, not a copy.
func analyticsFixture(t *testing.T, configExtra string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n"+configExtra)
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}
{% include "analytics.html" %}`)
	writeFile(t, root, "components/analytics.html", readScaffold(t, "scaffold/components/analytics.html"))
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	return root
}

func TestAnalyticsEmitsNothingWhenUnset(t *testing.T) {
	root := analyticsFixture(t, "")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if html := read(t, root, "index.html"); strings.Contains(html, "<script") {
		t.Errorf("analytics should emit nothing when unset:\n%s", html)
	}
}

func TestAnalyticsCloudflare(t *testing.T) {
	root := analyticsFixture(t, "analytics:\n  provider: cloudflare\n  token: abc123\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	html := read(t, root, "index.html")
	if !strings.Contains(html, `data-cf-beacon='{"token": "abc123"}'`) {
		t.Errorf("cloudflare analytics not emitted:\n%s", html)
	}
}

func TestAnalyticsPlausible(t *testing.T) {
	root := analyticsFixture(t, "analytics:\n  provider: plausible\n  domain: example.com\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	html := read(t, root, "index.html")
	if !strings.Contains(html, `data-domain="example.com"`) || !strings.Contains(html, "plausible.io/js/script.js") {
		t.Errorf("plausible analytics not emitted:\n%s", html)
	}
}

func TestAnalyticsFathom(t *testing.T) {
	root := analyticsFixture(t, "analytics:\n  provider: fathom\n  site_id: XYZ987\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	html := read(t, root, "index.html")
	if !strings.Contains(html, `data-site="XYZ987"`) || !strings.Contains(html, "usefathom.com/script.js") {
		t.Errorf("fathom analytics not emitted:\n%s", html)
	}
}

// TestJSONLDValidates builds a project with the real scaffold business.yaml
// and jsonld.html, then parses the emitted <script type="application/ld+json">
// block as JSON — guarding against a stray/missing comma in the conditional
// fields breaking the output.
func TestJSONLDValidates(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "data/business.yaml", readScaffold(t, "scaffold/data/business.yaml"))
	writeFile(t, root, "components/jsonld.html", readScaffold(t, "scaffold/components/jsonld.html"))
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}
{% include "jsonld.html" %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	parsed := parseJSONLD(t, read(t, root, "index.html"))
	if parsed["name"] != "My Business" {
		t.Errorf("jsonld name = %v, want %q", parsed["name"], "My Business")
	}
}

// parseJSONLD extracts the <script type="application/ld+json"> block from
// rendered HTML and parses it as JSON.
func parseJSONLD(t *testing.T, html string) map[string]any {
	t.Helper()
	start := strings.Index(html, "<script type=\"application/ld+json\">")
	if start == -1 {
		t.Fatalf("jsonld script not emitted:\n%s", html)
	}
	start += len("<script type=\"application/ld+json\">")
	end := strings.Index(html[start:], "</script>")
	if end == -1 {
		t.Fatalf("jsonld script not closed:\n%s", html)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(html[start:start+end]), &parsed); err != nil {
		t.Fatalf("jsonld did not parse as JSON: %v\n%s", err, html[start:start+end])
	}
	return parsed
}

// TestJSONLDEscapesQuotes guards the common case that broke under HTML
// autoescape: a business name with an apostrophe (Joe's Pizza) must
// round-trip through the emitted JSON verbatim, not as Joe&#39;s.
func TestJSONLDEscapesQuotes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "data/business.yaml", `name: Joe's "Best" Cameras
description: LA's top shop <est. 1999>
sameAs:
  - https://example.com/?a=1&b=2
`)
	writeFile(t, root, "components/jsonld.html", readScaffold(t, "scaffold/components/jsonld.html"))
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}
{% include "jsonld.html" %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\n---\nhome\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	parsed := parseJSONLD(t, read(t, root, "index.html"))
	if parsed["name"] != `Joe's "Best" Cameras` {
		t.Errorf("jsonld name = %v, want the raw quoted name", parsed["name"])
	}
	if parsed["description"] != "LA's top shop <est. 1999>" {
		t.Errorf("jsonld description = %v", parsed["description"])
	}
	sameAs, _ := parsed["sameAs"].([]any)
	if len(sameAs) != 1 || sameAs[0] != "https://example.com/?a=1&b=2" {
		t.Errorf("jsonld sameAs = %v", parsed["sameAs"])
	}
}

// TestShopifyBuyRenders exercises the real scaffold shopify-buy.html snippet
// with a store domain/token from config.yaml and a product id from page
// frontmatter — the two config slots the usage note documents.
func TestShopifyBuyRenders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", `site_name: Fixture
base_url: https://fixture.test
shopify:
  domain: fixture-store.myshopify.com
  storefront_access_token: fixture-token
`)
	writeFile(t, root, "components/shopify-buy.html", readScaffold(t, "scaffold/components/shopify-buy.html"))
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}
{% include "shopify-buy.html" %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\nshopify_product_id: \"123456789\"\n---\nhome\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	html := read(t, root, "index.html")
	for _, want := range []string{
		`<div id="product-component-123456789">`,
		`domain: "fixture-store.myshopify.com"`,
		`storefrontAccessToken: "fixture-token"`,
		`id: "123456789"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("shopify-buy.html missing %q:\n%s", want, html)
		}
	}
}

// TestStripeBuyRenders exercises the real scaffold stripe-buy.html snippet
// with a publishable key from config.yaml and a buy button id from page
// frontmatter — the two config slots the usage note documents.
func TestStripeBuyRenders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", `site_name: Fixture
base_url: https://fixture.test
stripe:
  publishable_key: pk_live_fixture
`)
	writeFile(t, root, "components/stripe-buy.html", readScaffold(t, "scaffold/components/stripe-buy.html"))
	writeFile(t, root, "layouts/base.html", `{% block content %}{{ content | safe }}{% endblock %}
{% include "stripe-buy.html" %}`)
	writeFile(t, root, "content/index.md", "---\nlayout: base.html\nstripe_buy_button_id: buy_btn_fixture\n---\nhome\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	html := read(t, root, "index.html")
	for _, want := range []string{
		`<script async src="https://js.stripe.com/v3/buy-button.js"></script>`,
		`buy-button-id="buy_btn_fixture"`,
		`publishable-key="pk_live_fixture"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("stripe-buy.html missing %q:\n%s", want, html)
		}
	}
}

// TestHeadPartialRenders exercises the real scaffold head.html against the
// acceptance criteria: correct title, description, canonical, and OG +
// Twitter tags from frontmatter/config.
func TestHeadPartialRenders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nbase_url: https://fixture.test\n")
	writeFile(t, root, "components/head.html", readScaffold(t, "scaffold/components/head.html"))
	writeFile(t, root, "components/jsonld.html", readScaffold(t, "scaffold/components/jsonld.html"))
	writeFile(t, root, "components/favicons.html", readScaffold(t, "scaffold/components/favicons.html"))
	writeFile(t, root, "layouts/base.html", `<head>{% include "head.html" %}</head>{% block content %}{{ content | safe }}{% endblock %}`)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\ndescription: A test page.\nimage: /og.jpg\nlayout: base.html\n---\nbody\n")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}

	html := read(t, root, "index.html")
	for _, want := range []string{
		"<title>Home — Fixture</title>",
		`<meta name="description" content="A test page.">`,
		`<link rel="canonical" href="https://fixture.test/">`,
		`<meta property="og:title" content="Home">`,
		`<meta property="og:description" content="A test page.">`,
		`<meta property="og:image" content="https://fixture.test/og.jpg">`,
		`<meta name="twitter:card" content="summary_large_image">`,
		`<meta name="twitter:image" content="https://fixture.test/og.jpg">`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("head.html missing %q:\n%s", want, html)
		}
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

// TestBuildSiteRejectsDangerousOutputDir guards the os.RemoveAll in buildSite:
// an output_dir pointing at the project root, a source dir, or anywhere
// outside the project must fail before anything is deleted.
func TestBuildSiteRejectsDangerousOutputDir(t *testing.T) {
	for _, dir := range []string{".", "..", "../elsewhere", "content", "content/out", "layouts", "components", "data", "static"} {
		root := t.TempDir()
		writeFile(t, root, "config.yaml", "site_name: X\noutput_dir: "+dir+"\n")
		writeFile(t, root, "content/index.md", "home\n")
		if _, err := buildSite(root); err == nil {
			t.Errorf("output_dir %q should be rejected", dir)
		}
		if _, statErr := os.Stat(filepath.Join(root, "content", "index.md")); statErr != nil {
			t.Errorf("output_dir %q: source content was deleted", dir)
		}
	}

	// a harmless custom dir still works
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: X\noutput_dir: dist\n")
	writeFile(t, root, "content/index.md", "home\n")
	if _, err := buildSite(root); err != nil {
		t.Errorf("output_dir \"dist\" should build: %v", err)
	}
}

// TestBuildSiteRejectsEscapingPermalink guards the page-write path the same
// way checkOutputDir guards the wipe: a permalink must never resolve outside
// the output dir.
func TestBuildSiteRejectsEscapingPermalink(t *testing.T) {
	for _, permalink := range []string{"../evil", "/../evil", "/a/../../evil"} {
		root := t.TempDir()
		writeFile(t, root, "config.yaml", "site_name: X\n")
		writeFile(t, root, "content/evil.md", "---\npermalink: "+permalink+"\n---\nowned\n")
		if _, err := buildSite(root); err == nil {
			t.Errorf("permalink %q should be rejected", permalink)
		}
		if _, statErr := os.Stat(filepath.Join(root, "evil.html")); statErr == nil {
			t.Errorf("permalink %q: wrote outside the output dir", permalink)
		}
	}
}

func TestBuildSiteNoConfig(t *testing.T) {
	_, err := buildSite(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "config.yaml") {
		t.Errorf("err = %v", err)
	}
}
