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

// scaffoldStyles is the seed stylesheet copied to static/style.css. Vanilla CSS,
// Safaii conventions: @layer config/reset/elements/components, oklch tokens as
// custom properties, rem everywhere (px only for border/outline width), fluid
// type via clamp() with cqi, and enhancements gated behind user preferences.
// Tokens mirror DESIGN.md one-to-one (colors.ink -> --color-ink, etc.).
const scaffoldStyles = `@layer config, reset, elements, components;

@layer config {
  :root {
    --color-ink: oklch(20% 0.01 250);
    --color-paper: oklch(98% 0.005 95);
    --color-accent: oklch(55% 0.18 25);
    --color-accent-strong: oklch(48% 0.18 25);
    --color-muted: oklch(55% 0.01 250);
    --color-hairline: oklch(85% 0.01 250);

    --space-xs: 0.5rem;
    --space-sm: 0.75rem;
    --space-md: 1rem;
    --space-lg: 1.5rem;
    --space-xl: 2.5rem;

    --radius-sm: 0.25rem;
    --radius-md: 0.5rem;
    --radius-lg: 1rem;

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
    max-width: 70rem;
    margin-inline: auto;
    padding: var(--space-xl) var(--space-lg);
    background: var(--color-paper);
    color: var(--color-ink);
    font-family: var(--font-serif);
    font-size: 1rem;
    line-height: 1.6;
  }

  h1,
  h2,
  h3 {
    margin-block: var(--space-lg) var(--space-md);
    font-family: var(--font-sans);
    line-height: 1.1;
  }

  /* Fluid headings: rem minimum, cqi ideal, rem max. */
  h1 { font-size: clamp(2rem, 6cqi, 3.5rem); }
  h2 { font-size: clamp(1.5rem, 4cqi, 2.25rem); }
  h3 { font-size: 1.25rem; }

  p { margin-block: var(--space-md); }

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
    padding: var(--space-sm) var(--space-lg);
    border-radius: var(--radius-md);
    background: var(--color-accent);
    color: var(--color-paper);
    font-family: var(--font-sans);
    text-decoration: none;
  }

  @media (prefers-reduced-motion: no-preference) {
    .button { transition: background-color 150ms ease; }
  }

  @media (hover: hover) {
    .button:hover { background: var(--color-accent-strong); }
  }

  /* Component with descendants: scoped so its selectors never leak out. */
  @scope (.card) {
    :scope {
      padding: var(--space-lg);
      border: 1px solid var(--color-hairline);
      border-radius: var(--radius-lg);
      background: var(--color-paper);
    }

    @media (prefers-reduced-transparency: no-preference) {
      :scope { background: oklch(98% 0.005 95 / 0.75); }
    }

    h3 { margin-block: 0 var(--space-sm); }

    p { color: var(--color-muted); }
  }
}
`

// scaffoldDesignMD is the seed DESIGN.md: YAML front matter (machine-readable
// tokens, mapped 1:1 to the CSS custom properties in scaffoldStyles) plus prose
// rationale. fontSize stores the rem MINIMUM of each fluid step; the clamp()
// ideal/max live in the Typography prose so the front matter stays lint-clean.
// Keep token values identical to scaffoldStyles. No backticks in the body.
const scaffoldDesignMD = `---
name: My Site
version: 1.0.0
description: Seed design system for a Garp site. Tokens map one-to-one to CSS custom properties in static/style.css.
colors:
  ink: oklch(20% 0.01 250)
  paper: oklch(98% 0.005 95)
  accent: oklch(55% 0.18 25)
  muted: oklch(55% 0.01 250)
  hairline: oklch(85% 0.01 250)
typography:
  body:
    fontFamily: Georgia, Cambria, "Times New Roman", Times, serif
    fontSize: 1rem
    lineHeight: 1.6
  heading:
    fontFamily: system-ui, sans-serif
    fontSize: 1.5rem
    lineHeight: 1.1
  display:
    fontFamily: system-ui, sans-serif
    fontSize: 2rem
    lineHeight: 1.1
rounded:
  sm: 0.25rem
  md: 0.5rem
  lg: 1rem
spacing:
  xs: 0.5rem
  sm: 0.75rem
  md: 1rem
  lg: 1.5rem
  xl: 2.5rem
components:
  button:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.paper}"
    rounded: "{rounded.md}"
    padding: "{spacing.sm}"
  buttonHover:
    backgroundColor: oklch(48% 0.18 25)
    textColor: "{colors.paper}"
    rounded: "{rounded.md}"
    padding: "{spacing.sm}"
  card:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.lg}"
    padding: "{spacing.lg}"
---

# Overview

This is the seed design system for a Garp site. The front matter above is the
single source of truth for design tokens; every token maps one-to-one to a CSS
custom property in static/style.css (colors.ink to --color-ink, spacing.md to
--space-md, rounded.lg to --radius-lg). Edit tokens here and mirror the change
in the stylesheet, or run npx design.md to lint and diff.

The stylesheet uses cascade layers in a fixed order: config (the tokens as
custom properties), reset, elements, then components.

# Colors

Colors are authored in oklch for perceptual uniformity. ink is the near-black
body text, paper the warm off-white background, accent the warm red used for
links and the primary button, muted a grey for secondary text, and hairline a
light grey for borders. A darker accent variant is used for button hover.

# Typography

Two families: a serif for body copy and a sans (system-ui) for headings and UI.
Body text is a fixed 1rem. Headings are fluid via clamp(), which the front
matter cannot express in a single dimension, so fontSize records only the rem
MINIMUM. The full curves are:

- heading: clamp(1.5rem, 4cqi, 2.25rem)
- display: clamp(2rem, 6cqi, 3.5rem)

Fluid type resolves against the container, so body sets container-type:
inline-size for the cqi unit to work.

# Layout

The body is centered with a max-width measure and padded with the spacing
scale (xs 0.5rem through xl 2.5rem). All sizing uses relative units; the only
pixel values are border-width and outline-width, kept crisp at any zoom.

# Elevation and Depth

Flat by default. Separation comes from the hairline border and the paper/ink
contrast rather than shadows.

# Shapes

Three corner radii: sm 0.25rem for small controls, md 0.5rem for buttons, lg
1rem for cards.

# Components

- button: a root-only component (a single plain CSS rule). Hover is gated
  behind @media (hover: hover); its color transition is gated behind
  prefers-reduced-motion.
- card: a bordered container styled with @scope so its descendant rules never
  leak. Any translucency is gated behind prefers-reduced-transparency.

# Dos and Donts

- Do reference tokens with var() in CSS and with brace paths in this document.
- Do gate hover, motion, and translucency behind the matching media query.
- Do keep relative units; reserve px for border-width and outline-width.
- Dont hard-code color or spacing values that a token already covers.
- Dont let the front matter and the stylesheet drift apart.
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
		"DESIGN.md":         scaffoldDesignMD,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			return err
		}
	}

	fmt.Printf("Created %s\n\n  cd %s\n  garp serve\n", root, root)
	return nil
}
