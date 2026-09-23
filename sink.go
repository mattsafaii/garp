package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The kitchen sink is a generated design-system reference page: the
// project's tokens (custom properties declared in the stylesheet's
// @layer tokens block) rendered as swatches and specimens, plus a fixed
// set of semantic elements exercised by the project's own CSS. Token
// *values* are never computed here — previews reference var(--token) and
// the browser resolves them against the live stylesheet, so the page
// cannot drift. It exists in two tiers: /_garp/sink on the dev server
// (regenerated per request), and an opt-in build synthesis (sink: true
// in config.yaml) writing kitchen-sink.html, author-overridable the
// same way as sitemap.xml.

// token is one custom property declared in @layer tokens, with its raw
// authored value (comments stripped) for display.
type token struct {
	Name  string
	Value string
}

// tokenGroup is a titled set of tokens sharing a preview treatment.
type tokenGroup struct {
	Title  string
	Tokens []token
}

// tokenRe deliberately restricts names to [a-zA-Z0-9-]: names are
// embedded in style attributes, so the character set is the injection
// guard.
var tokenRe = regexp.MustCompile(`(--[a-zA-Z0-9-]+)\s*:\s*([^;{}]*);`)

var cssCommentRe = regexp.MustCompile(`/\*.*?\*/`)

// extractTokens returns the custom properties declared inside the
// stylesheet's @layer tokens block, in declaration order. The tokens
// layer is the contract (per the studio CSS conventions); a stylesheet
// without one yields no tokens.
func extractTokens(css string) []token {
	body, ok := layerBlock(css, "tokens")
	if !ok {
		return nil
	}
	var tokens []token
	for _, m := range tokenRe.FindAllStringSubmatch(body, -1) {
		value := strings.TrimSpace(cssCommentRe.ReplaceAllString(m[2], ""))
		tokens = append(tokens, token{Name: m[1], Value: value})
	}
	return tokens
}

// layerBlock returns the brace-matched body of `@layer <name> { ... }`.
// The layer *statement* form (`@layer reset, tokens, ...;`) has no body
// and is skipped naturally because it has no `{`.
func layerBlock(css, name string) (string, bool) {
	re := regexp.MustCompile(`@layer\s+` + name + `\s*\{`)
	loc := re.FindStringIndex(css)
	if loc == nil {
		return "", false
	}
	depth := 1
	start := loc[1]
	for i := start; i < len(css); i++ {
		switch css[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return css[start:i], true
			}
		}
	}
	return "", false
}

// groupTokens buckets tokens by name prefix into display groups, keeping
// declaration order within each. Unrecognized tokens land in Other and
// get no preview — name and value are still worth listing.
func groupTokens(tokens []token) []tokenGroup {
	titles := []string{"Colors", "Spacing", "Typography", "Radius", "Shadows", "Other"}
	buckets := make(map[string][]token, len(titles))
	for _, t := range tokens {
		buckets[groupOf(t.Name)] = append(buckets[groupOf(t.Name)], t)
	}
	var groups []tokenGroup
	for _, title := range titles {
		if len(buckets[title]) > 0 {
			groups = append(groups, tokenGroup{Title: title, Tokens: buckets[title]})
		}
	}
	return groups
}

func groupOf(name string) string {
	switch {
	case strings.HasPrefix(name, "--color-"):
		return "Colors"
	case strings.HasPrefix(name, "--space-"), name == "--gutter", name == "--stack":
		return "Spacing"
	case strings.HasPrefix(name, "--font-"), strings.HasPrefix(name, "--text-"),
		strings.HasPrefix(name, "--line-height"), strings.HasPrefix(name, "--measure"),
		strings.HasPrefix(name, "--letter-"):
		return "Typography"
	case strings.HasPrefix(name, "--radius"):
		return "Radius"
	case strings.HasPrefix(name, "--shadow"):
		return "Shadows"
	}
	return "Other"
}

// previewStyle maps a token to the inline style of its preview cell —
// always a var() reference, never a computed value, so the browser
// resolves it against the live stylesheet. An empty string means no
// visual preview makes sense for this token.
func previewStyle(name string) string {
	switch {
	case strings.HasPrefix(name, "--color-"):
		return "background-color: var(" + name + ")"
	case strings.HasPrefix(name, "--space-"), name == "--gutter", name == "--stack":
		return "width: var(" + name + ")"
	case strings.HasPrefix(name, "--radius"):
		return "border-radius: var(" + name + ")"
	case strings.HasPrefix(name, "--shadow"):
		return "box-shadow: var(" + name + ")"
	case strings.HasPrefix(name, "--font-"):
		return "font-family: var(" + name + ")"
	case strings.HasPrefix(name, "--text-"):
		return "font-size: var(" + name + ")"
	}
	return ""
}

// previewText is the specimen content inside a preview cell, where one
// applies.
func previewText(name string) string {
	if strings.HasPrefix(name, "--font-") || strings.HasPrefix(name, "--text-") {
		return "Aa Bb Cc 0123"
	}
	return ""
}

// renderSinkPage produces the standalone kitchen-sink HTML page. It
// links the project's own /style.css; its own chrome styles are
// namespaced under .gsink- so they can't collide with project CSS.
func renderSinkPage(groups []tokenGroup, siteName string) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	b.WriteString("<meta name=\"robots\" content=\"noindex\">\n")
	b.WriteString("<title>Kitchen sink — " + html.EscapeString(siteName) + "</title>\n")
	b.WriteString("<link rel=\"stylesheet\" href=\"/style.css\">\n")
	b.WriteString("<style>\n" + sinkChrome + "</style>\n</head>\n<body>\n<main>\n")
	b.WriteString("<h1>Kitchen sink</h1>\n")
	b.WriteString("<p>Every token from <code>static/style.css</code>'s <code>@layer tokens</code>, and the standard elements, rendered by this site's own stylesheet. Generated by garp — token values are live <code>var()</code> references, so this page is always current.</p>\n")

	if len(groups) == 0 {
		b.WriteString("<p><em>No tokens found — static/style.css has no <code>@layer tokens</code> block.</em></p>\n")
	}
	for _, g := range groups {
		b.WriteString("<h2>" + g.Title + "</h2>\n<ul class=\"gsink-tokens\">\n")
		for _, t := range g.Tokens {
			b.WriteString("<li>")
			if style := previewStyle(t.Name); style != "" {
				b.WriteString("<span class=\"gsink-preview gsink-" + strings.ToLower(g.Title) + "\" style=\"" + style + "\">" + previewText(t.Name) + "</span>")
			}
			b.WriteString("<code class=\"gsink-name\">" + t.Name + "</code>")
			b.WriteString("<code class=\"gsink-value\">" + html.EscapeString(t.Value) + "</code>")
			b.WriteString("</li>\n")
		}
		b.WriteString("</ul>\n")
	}

	b.WriteString(sinkElements)
	b.WriteString("</main>\n</body>\n</html>\n")
	return b.String()
}

// serveSink answers /_garp/sink on the dev server: parse the source
// stylesheet fresh on every request so the page is current even before
// the next rebuild copies it to the output dir.
func serveSink(root string, w http.ResponseWriter) {
	css, err := os.ReadFile(filepath.Join(root, "static", "style.css"))
	if err != nil {
		http.Error(w, "static/style.css not found", http.StatusNotFound)
		return
	}
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, renderSinkPage(groupTokens(extractTokens(string(css))), cfg.SiteName))
}

// synthesizeSink writes kitchen-sink.html into outDir when the author
// opts in with `sink: true` in config.yaml. Author-overridable exactly
// like sitemap.xml: a static/kitchen-sink.html wins and synthesis is
// skipped. Returns the number of files written (0 or 1).
func synthesizeSink(root, outDir string, cfg *Config) (int, error) {
	if opt, _ := cfg.Site["sink"].(bool); !opt {
		return 0, nil
	}
	if _, err := os.Stat(filepath.Join(root, "static", "kitchen-sink.html")); err == nil {
		return 0, nil
	} else if !os.IsNotExist(err) {
		return 0, err
	}
	css, err := os.ReadFile(filepath.Join(root, "static", "style.css"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "garp: sink: true but static/style.css not found — skipping kitchen-sink.html")
		return 0, nil
	}
	page := renderSinkPage(groupTokens(extractTokens(string(css))), cfg.SiteName)
	if err := os.WriteFile(filepath.Join(outDir, "kitchen-sink.html"), []byte(page), 0o644); err != nil {
		return 0, err
	}
	return 1, nil
}

// sinkChrome is the page's own minimal styling, namespaced so project
// CSS and sink chrome can't touch each other.
const sinkChrome = `.gsink-tokens { list-style: none; padding: 0; }
.gsink-tokens > li { display: flex; align-items: center; gap: 0.75rem; padding-block: 0.25rem; }
.gsink-preview { flex: none; display: inline-block; min-width: 2rem; min-height: 2rem; border: 1px solid rgba(128,128,128,0.35); }
.gsink-spacing { min-width: 0; height: 1rem; background: rgba(128,128,128,0.5); border: 0; }
.gsink-typography { border: 0; min-height: 0; white-space: nowrap; }
.gsink-name { font-weight: 600; }
.gsink-value { opacity: 0.6; font-size: 0.8em; overflow-wrap: anywhere; }
`

// sinkElements is the fixed semantic-HTML section — markup the
// project's stylesheet is expected to style. It needs no generation:
// the browser applies the live CSS, so it is current by construction.
const sinkElements = `<h2>Elements</h2>

<h1>Heading level 1</h1>
<h2>Heading level 2</h2>
<h3>Heading level 3</h3>
<h4>Heading level 4</h4>

<p>A paragraph with <a href="#">a link</a>, <strong>strong text</strong>, <em>emphasis</em>, <code>inline code</code>, <mark>marked text</mark>, <small>small text</small>, an <abbr title="abbreviation">abbr</abbr>, <sub>sub</sub> and <sup>sup</sup>, and <s>struck text</s>. Sphinx of black quartz, judge my vow.</p>

<blockquote>
<p>Those people who think they know everything are a great annoyance to those of us who do.</p>
<cite>Isaac Asimov</cite>
</blockquote>

<ul>
<li>Unordered list item</li>
<li>Another item
<ul><li>Nested item</li><li>Nested item</li></ul>
</li>
<li>A third item</li>
</ul>

<ol>
<li>Ordered list item</li>
<li>Another item</li>
<li>A third item</li>
</ol>

<dl>
<dt>Definition term</dt>
<dd>The definition's description text.</dd>
<dt>Another term</dt>
<dd>Another description.</dd>
</dl>

<pre><code>&lt;section class="example"&gt;
  &lt;p&gt;Block code sample&lt;/p&gt;
&lt;/section&gt;</code></pre>

<table>
<thead><tr><th>Header one</th><th>Header two</th><th>Header three</th></tr></thead>
<tbody>
<tr><td>Cell</td><td>Cell</td><td>Cell</td></tr>
<tr><td>Cell</td><td>Cell</td><td>Cell</td></tr>
</tbody>
</table>

<hr>

<details>
<summary>A details disclosure</summary>
<p>Content revealed on toggle.</p>
</details>

<figure>
<img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='320' height='180'%3E%3Crect width='320' height='180' fill='%23888'/%3E%3C/svg%3E" alt="Placeholder image" width="320" height="180">
<figcaption>A figure with a caption.</figcaption>
</figure>

<form action="#" method="post">
<p>
<label for="gsink-text">Text input</label>
<input id="gsink-text" name="text" type="text" placeholder="Placeholder">
</p>
<p>
<label for="gsink-email">Email input</label>
<input id="gsink-email" name="email" type="email" placeholder="you@example.com" required>
</p>
<p>
<label for="gsink-select">Select</label>
<select id="gsink-select" name="select">
<option>Option one</option>
<option>Option two</option>
<optgroup label="Group">
<option>Grouped option</option>
</optgroup>
</select>
</p>
<p>
<label for="gsink-textarea">Textarea</label>
<textarea id="gsink-textarea" name="textarea" rows="3" placeholder="Multi-line text"></textarea>
</p>
<fieldset>
<legend>Radios and checkboxes</legend>
<label><input type="radio" name="gsink-radio" checked> Radio one</label>
<label><input type="radio" name="gsink-radio"> Radio two</label>
<label><input type="checkbox" checked> Checkbox</label>
</fieldset>
<p>
<label for="gsink-range">Range</label>
<input id="gsink-range" name="range" type="range">
</p>
<p>
<button type="submit">Submit button</button>
<button type="button" disabled>Disabled button</button>
</p>
</form>
`
