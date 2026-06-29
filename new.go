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

// scaffoldStyles is the seed stylesheet copied to static/style.css. Vanilla CSS
// derived from the safaii-css starter (the canonical source of truth): @layer
// config/reset/elements/components, oklch tokens as custom properties, the
// --gutter/--stack/--measure spacing tokens (other sizes via calc()), a single
// --radius, fluid type via clamp() with cqi, and enhancements gated behind user
// preferences. The stylesheet is the design spec; keep it in sync with the
// safaii-css references/starter.css.
const scaffoldStyles = `@layer config, reset, elements, components;

@layer config {
  :root {
    --gutter: clamp(1ch, 2.5vmax, 3ch); /* inline spacing */
    --stack: clamp(1.25ex, 2.5vmax, 1.75ex); /* block spacing */
    --measure: 64ch; /* line length */
    --radius: 0.25rem; /* one rem token, not literal px */

    --color-bg: oklch(98% 0.005 95);
    --color-text: oklch(20% 0.01 250);
    --color-accent: oklch(55% 0.18 25);
    --color-border: oklch(85% 0.01 250);

    --font-serif: Georgia, Cambria, "Times New Roman", Times, serif;
    --font-sans: system-ui, sans-serif;
  }
}

@layer reset {
  *,
  *::before,
  *::after { box-sizing: border-box; }

  * { margin: 0; }

  img,
  picture,
  svg { display: block; max-width: 100%; }

  a { color: inherit; }
}

@layer elements {
  body {
    container-type: inline-size;
    max-inline-size: var(--measure);
    margin-inline: auto;
    padding-block: calc(var(--stack) * 3);
    padding-inline: var(--gutter);
    background: var(--color-bg);
    color: var(--color-text);
    font-family: var(--font-serif);
    font-size: 1rem;
    line-height: 1.6;
  }

  h1,
  h2,
  h3 {
    margin-block: calc(var(--stack) * 2) var(--stack);
    font-family: var(--font-sans);
    line-height: 1.1;
  }

  /* Fluid headings: rem minimum, cqi ideal, rem max. */
  h1 { font-size: clamp(2rem, 6cqi, 3.5rem); }
  h2 { font-size: clamp(1.5rem, 4cqi, 2.25rem); }
  h3 { font-size: 1.25rem; }

  p { margin-block: var(--stack); }

  a { color: var(--color-accent); }

  :focus-visible {
    outline: 2px solid var(--color-accent);
    outline-offset: 0.125rem;
  }
}

@layer components {
  /* Root-only component: a single plain rule, no descendants. */
  .button {
    display: inline-block;
    padding: var(--stack) var(--gutter);
    border-radius: var(--radius);
    background: var(--color-accent);
    color: var(--color-bg);
    font-family: var(--font-sans);
    text-decoration: none;
  }

  @media (prefers-reduced-motion: no-preference) {
    .button { transition: background-color 150ms ease; }
  }

  @media (hover: hover) {
    /* Darker accent derived inline, not a separate token. */
    .button:hover { background: oklch(from var(--color-accent) calc(l - 0.07) c h); }
  }

  /* Component with descendants: scoped so its selectors never leak out. */
  @scope (.card) {
    :scope {
      padding: var(--gutter);
      border: 1px solid var(--color-border);
      border-radius: var(--radius);
      background: var(--color-bg);
    }

    @media (prefers-reduced-transparency: no-preference) {
      :scope { background: oklch(from var(--color-bg) l c h / 0.75); }
    }

    h3 { margin-block: 0 var(--stack); }

    /* Muted text derived from --color-text, not a separate token. */
    p { color: color-mix(in oklch, var(--color-text), var(--color-bg) 40%); }
  }
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
