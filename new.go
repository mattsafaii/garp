package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const scaffoldConfig = `site_name: My Site
base_url: https://example.com
`

const scaffoldLayout = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{% if title %}{{ title }} — {% endif %}{{ site.site_name }}</title>
  <link rel="stylesheet" href="/style.css">
</head>
<body>
  {% block content %}{{ content | safe }}{% endblock %}
</body>
</html>
`

const scaffoldIndex = `---
title: Home
layout: base.html
---

# Welcome to Garp

This page is ` + "`content/index.md`" + `, rendered through ` + "`layouts/base.html`" + `.
Edit it, run ` + "`garp serve`" + `, and refresh.
`

// scaffoldStyles is the seed stylesheet copied to static/style.css. It is a
// verbatim copy of the safaii-css starter (the canonical source of truth):
// ~/.claude/skills/safaii-css/references/starter.css. Do not hand-edit this
// const — when the conventions move, refresh it from starter.css so `garp new`
// keeps emitting exactly the canonical starter.
const scaffoldStyles = `/* Starter stylesheet — extracted from mattsafaii.com (src/css/style.css).
   Copy as the project's single stylesheet. Fill in the palette, keep the bones. */

/* ── Layers ──────────────────────────────────────────────── */

@layer config, reset, elements, components;

@layer config {
	:root {
		color-scheme: light;
		font-size: clamp(95%, 85% + 0.5dvi, 115%);
		--gutter: clamp(1ch, 2.5vmax, 3ch); /* inline spacing */
		--stack: clamp(1.25ex, 2.5vmax, 1.75ex); /* block spacing */
		--measure: 64ch;
		--radius: 0.25rem; /* corner radius — one rem token, not literal px */
		/* Project palette — always oklch */
		--color-bg: oklch(98% 0 0);
		--color-text: oklch(20% 0 0);
		--color-accent: oklch(45% 0.15 25);
		accent-color: var(--color-accent);
	}

	@view-transition {
		navigation: auto;
	}
}

@layer reset {
	*,
	*::before,
	*::after {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
		font-kerning: normal;
	}

	input,
	button,
	textarea,
	select {
		font: inherit;
	}

	@media (forced-colors: active) {
		:where(button) {
			border: 1px solid;
		}
	}
}

@layer elements {
	@media (prefers-reduced-motion: no-preference) {
		:root {
			scroll-behavior: smooth;
		}
	}

	body {
		-webkit-font-smoothing: antialiased;
		font: 1rem / 1.35 Georgia, serif; /* project type goes here */
		color: var(--color-text);
		background-color: var(--color-bg);
		display: grid;
		grid-template-columns:
			[bleed-start] minmax(var(--gutter), 1fr)
			[content-start] minmax(0, var(--measure))
			[content-end] minmax(var(--gutter), 1fr)
			[bleed-end];
		align-items: start;
		padding-block: 8vh calc(var(--gutter) * 2);
		overflow-x: clip;
	}

	main {
		grid-column: bleed;
		display: grid;
		grid-template-columns: subgrid;
		align-items: start;
	}

	main > * {
		grid-column: content;
	}

	main > * + * {
		margin-block-start: calc(var(--stack) * 2);
	}

	a {
		color: inherit;
		text-underline-offset: 0.15em;
		text-decoration-thickness: 0.05em;
	}

	a:hover {
		text-decoration-color: var(--color-accent);
	}

	h1,
	h2,
	h3 {
		margin-block-end: calc(var(--stack) / 2);
		text-wrap: balance;
		overflow-wrap: break-word;
		hyphens: auto;
	}

	p,
	ul,
	ol {
		margin-block-end: var(--stack);
		text-wrap: pretty;
	}

	ul,
	ol {
		padding-inline-start: 2ch;
	}

	::selection {
		background: var(--color-accent);
		color: var(--color-bg);
	}

	:where(:focus-visible) {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	:where(img, svg, video, iframe) {
		max-inline-size: 100%;
		block-size: auto;
	}

	:where(svg) {
		fill: currentColor;
	}

	article {
		content-visibility: auto;
	}
}

@layer components {
	/* Each component: @scope to the root; use ` + "`to (...)`" + ` to stop at
	   nested component boundaries. Example shape:

	@scope (.site-nav) {
		:scope {
			display: flex;
			gap: var(--gutter);
		}

		a {
			text-decoration: none;
		}
	}

	@scope (.card) to (.card-body) {
		img {
			border-radius: var(--radius);
		}
	}
	*/
}
`

func cmdNew(args []string) error {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp new <path>")
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	root := fs.Arg(0)

	if _, err := os.Stat(root); err == nil {
		return fmt.Errorf("%s already exists", root)
	}

	for _, dir := range []string{"content", "layouts", "components", "data", "static", "site"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"config.yaml":       scaffoldConfig,
		"layouts/base.html": scaffoldLayout,
		"content/index.md":  scaffoldIndex,
		"static/style.css":  scaffoldStyles,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			return err
		}
	}

	fmt.Printf("Created %s\n\n  cd %s\n  garp serve\n", root, root)
	return nil
}
