package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sinkFixtureCSS = `@layer reset, tokens, base, composition, block, utility, exception;

@layer tokens {
	:root {
		--color-bg: oklch(98% 0 0);
		--color-accent: oklch(55% 0.13 250); /* placeholder */
		--space-4: 1rem; /* 16 */
		--gutter: clamp(var(--space-4), 2.5vmax, 1.5rem);
		--font-sans: system-ui, sans-serif;
		--radius: 0.25rem;
		--custom-thing: 42;
	}

	@view-transition {
		navigation: auto;
	}
}

@layer base {
	:root {
		--not-a-token: red;
	}
}
`

// TestExtractTokens: only @layer tokens declarations count, order
// is preserved, comments are stripped from values, and nested blocks
// inside the layer don't break brace matching.
func TestExtractConfigTokens(t *testing.T) {
	tokens := extractTokens(sinkFixtureCSS)
	var names []string
	for _, tok := range tokens {
		names = append(names, tok.Name)
	}
	want := "--color-bg,--color-accent,--space-4,--gutter,--font-sans,--radius,--custom-thing"
	if got := strings.Join(names, ","); got != want {
		t.Errorf("token names = %s, want %s", got, want)
	}
	if tokens[1].Value != "oklch(55% 0.13 250)" {
		t.Errorf("comment not stripped from value: %q", tokens[1].Value)
	}
}

// TestExtractNoTokensLayer: a stylesheet without the tokens layer yields
// no tokens — the contract is explicit.
func TestExtractNoTokensLayer(t *testing.T) {
	if tokens := extractTokens(":root { --loose: 1; }"); tokens != nil {
		t.Errorf("expected no tokens, got %v", tokens)
	}
}

// TestGroupTokens: prefixes bucket correctly; --gutter counts as
// Spacing; unknown tokens land in Other.
func TestGroupTokens(t *testing.T) {
	groups := groupTokens(extractTokens(sinkFixtureCSS))
	got := map[string]int{}
	for _, g := range groups {
		got[g.Title] = len(g.Tokens)
	}
	for title, n := range map[string]int{"Colors": 2, "Spacing": 2, "Typography": 1, "Radius": 1, "Other": 1} {
		if got[title] != n {
			t.Errorf("group %s has %d tokens, want %d", title, got[title], n)
		}
	}
}

// TestRenderSinkPage: previews are var() references (never computed
// values), the page is noindex, and the elements section is present.
func TestRenderSinkPage(t *testing.T) {
	page := renderSinkPage(groupTokens(extractTokens(sinkFixtureCSS)), "Fixture")
	for _, want := range []string{
		`background-color: var(--color-accent)`,
		`width: var(--gutter)`,
		`<meta name="robots" content="noindex">`,
		`<link rel="stylesheet" href="/style.css">`,
		`<code class="gsink-name">--custom-thing</code>`,
		`<h2>Elements</h2>`,
		`<form action="#"`,
		`<blockquote>`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("sink page missing %q", want)
		}
	}
}

// TestDevSinkRoute: /_garp/sink regenerates from static/style.css per
// request.
func TestDevSinkRoute(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\n")
	writeFile(t, root, "static/style.css", sinkFixtureCSS)
	writeFile(t, root, "site/index.html", "<h1>home</h1>")
	rec := get(t, root, "/_garp/sink")
	if rec.Code != 200 {
		t.Fatalf("GET /_garp/sink = %d, want 200", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, "var(--color-bg)") {
		t.Errorf("sink route missing token preview:\n%s", body)
	}
}

// TestBuildSynthesizesSink: sink: true emits kitchen-sink.html; absent,
// nothing; a static/kitchen-sink.html wins over synthesis.
func TestBuildSynthesizesSink(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nsink: true\n")
	writeFile(t, root, "static/style.css", sinkFixtureCSS)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\n---\nhome\n")
	writeFile(t, root, "layouts/base.html", "{% block content %}{{ content | safe }}{% endblock %}")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	page := read(t, root, "kitchen-sink.html")
	if !strings.Contains(page, "var(--color-accent)") {
		t.Errorf("synthesized sink missing tokens:\n%s", page[:200])
	}
}

func TestBuildNoSinkByDefault(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\n")
	writeFile(t, root, "static/style.css", sinkFixtureCSS)
	writeFile(t, root, "content/index.md", "---\ntitle: Home\n---\nhome\n")
	writeFile(t, root, "layouts/base.html", "{% block content %}{{ content | safe }}{% endblock %}")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "site", "kitchen-sink.html")); !os.IsNotExist(err) {
		t.Error("kitchen-sink.html should not exist without sink: true")
	}
}

func TestBuildSinkAuthorOverride(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.yaml", "site_name: Fixture\nsink: true\n")
	writeFile(t, root, "static/style.css", sinkFixtureCSS)
	writeFile(t, root, "static/kitchen-sink.html", "<h1>authored sink</h1>")
	writeFile(t, root, "content/index.md", "---\ntitle: Home\n---\nhome\n")
	writeFile(t, root, "layouts/base.html", "{% block content %}{{ content | safe }}{% endblock %}")
	if _, err := buildSite(root); err != nil {
		t.Fatal(err)
	}
	if page := read(t, root, "kitchen-sink.html"); page != "<h1>authored sink</h1>" {
		t.Errorf("author's kitchen-sink.html should win, got %q", page)
	}
}

// TestScaffoldSinkConvention guards the contract between the scaffold's
// starter stylesheet and the sink parser: the starter must keep its
// tokens in @layer config where the sink finds them.
func TestScaffoldSinkConvention(t *testing.T) {
	css := readScaffold(t, "scaffold/static/style.css")
	tokens := extractTokens(css)
	if len(tokens) == 0 {
		t.Fatal("no tokens extracted from the scaffold starter — did @layer config change?")
	}
	names := map[string]bool{}
	for _, tok := range tokens {
		names[tok.Name] = true
	}
	for _, want := range []string{"--color-accent", "--space-4", "--radius"} {
		if !names[want] {
			t.Errorf("scaffold starter missing expected token %s", want)
		}
	}
}
