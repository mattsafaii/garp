package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// ciBinary is the platform Cloudflare Pages builds on — the binary that
// MUST be present for the committed repo to build in CI, even though Matt
// authors on darwin/arm64.
const ciBinary = "garp-linux-amd64"

func cmdHandoff(args []string) error {
	fs := flag.NewFlagSet("handoff", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp handoff")
	}
	fs.Parse(args)

	written, warning, err := runHandoff(".")
	if err != nil {
		return err
	}
	if warning != "" {
		fmt.Fprintln(os.Stderr, "garp: "+warning)
	}
	fmt.Printf("Wrote %d files\n", len(written))
	return nil
}

// runHandoff makes root self-sufficient: a bin/ of committed per-platform
// binaries plus a uname shim, and a generated HANDOFF.md describing this
// specific project. It's the bus-factor command — safe to re-run any time
// (it never touches content/layouts/data/static, only bin/ and HANDOFF.md),
// so unlike garp blog it doesn't refuse to overwrite. Returns the paths
// written and a non-fatal warning (e.g. the CI-platform binary wasn't found
// alongside the running one), if any.
func runHandoff(root string) (written []string, warning string, err error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, "", fmt.Errorf("config.yaml not found — is this a garp project?")
	}
	if err != nil {
		return nil, "", err
	}

	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return nil, "", err
	}

	binaries, err := collectBinaries()
	if err != nil {
		return nil, "", err
	}
	names := make([]string, 0, len(binaries))
	for name := range binaries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		data, err := os.ReadFile(binaries[name])
		if err != nil {
			return nil, "", fmt.Errorf("reading %s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(binDir, name), data, 0o755); err != nil {
			return nil, "", err
		}
		written = append(written, filepath.Join("bin", name))
	}
	if _, ok := binaries[ciBinary]; !ok {
		warning = fmt.Sprintf(
			"%s not found next to the running binary — Cloudflare Pages (linux/amd64) won't be able to build until it's added. Run `make release` in the garp repo and place bin/%s next to the garp binary you run `handoff` from, then rerun.",
			ciBinary, ciBinary,
		)
	}

	shim := "#!/bin/sh\nexec \"$(dirname \"$0\")/garp-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64\\|arm64/arm64/')\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(binDir, "garp"), []byte(shim), 0o755); err != nil {
		return nil, "", err
	}
	written = append(written, filepath.Join("bin", "garp"))

	doc, err := renderHandoffDoc(root, cfg)
	if err != nil {
		return nil, "", err
	}
	if err := os.WriteFile(filepath.Join(root, "HANDOFF.md"), []byte(doc), 0o644); err != nil {
		return nil, "", err
	}
	written = append(written, "HANDOFF.md")

	return written, warning, nil
}

// collectBinaries finds the running garp binary plus any platform siblings
// next to it (e.g. garp-linux-amd64, built by `make release` on Matt's
// Mac and left alongside the darwin/arm64 binary he actually runs). Keyed
// by the canonical bin/ filename, so the running binary and a same-named
// sibling collapse into one entry.
func collectBinaries() (map[string]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("locating the running binary: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	dir := filepath.Dir(exe)

	binaries := map[string]string{
		"garp-" + runtime.GOOS + "-" + runtime.GOARCH: exe,
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return binaries, nil // no sibling dir to scan is fine — current platform still works
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasPrefix(name, "garp-") {
			continue
		}
		if _, ok := binaries[name]; !ok {
			binaries[name] = filepath.Join(dir, name)
		}
	}
	return binaries, nil
}

// renderHandoffDoc builds HANDOFF.md from the project's actual shape — a
// garp-naive developer should be able to build, edit, and deploy this exact
// site from this file alone, per the bus-factor guarantee in CLAUDE.md.
func renderHandoffDoc(root string, cfg *Config) (string, error) {
	sections, err := topLevelContentDirs(filepath.Join(root, "content"))
	if err != nil {
		return "", err
	}
	dataFiles, err := topLevelFiles(filepath.Join(root, "data"))
	if err != nil {
		return "", err
	}

	hasBlog := dirExists(filepath.Join(root, "content", "blog"))
	hasFeed := fileExists(filepath.Join(root, "content", "feed.md"))
	hasFavicons := fileExists(filepath.Join(root, "static", "favicon.ico"))
	hasOG := dirExists(filepath.Join(root, "static", "og"))
	_, hasShopify := cfg.Site["shopify"]
	_, hasStripe := cfg.Site["stripe"]
	_, hasAnalytics := cfg.Site["analytics"]
	_, hasChat := cfg.Site["chat"]

	var b strings.Builder
	name := cfg.SiteName
	if name == "" {
		name = "this site"
	}
	fmt.Fprintf(&b, "# %s — Handoff\n\n", name)
	fmt.Fprintf(&b, "Generated by `garp handoff` on %s. This file, this repo, and the binaries in `bin/` are everything needed to build, edit, and deploy this site — no other tools or accounts required.\n\n", time.Now().Format("2006-01-02"))

	b.WriteString("## Quick start\n\n```sh\nbin/garp build   # writes the finished site to site/\nbin/garp serve   # build, serve locally, rebuild on save\n```\n\n")
	b.WriteString("`bin/garp` is a small shell shim that picks the right binary for your platform (macOS or Linux) automatically. If you're on another platform, install Go 1.26+ and run `go build -o bin/garp .` from this directory instead.\n\n")

	b.WriteString("## What's in this site\n\n")
	if cfg.BaseURL != "" {
		fmt.Fprintf(&b, "- **Live URL:** %s\n", cfg.BaseURL)
	}
	if len(sections) > 0 {
		fmt.Fprintf(&b, "- **Content sections:** %s\n", strings.Join(sections, ", "))
	}
	if len(dataFiles) > 0 {
		fmt.Fprintf(&b, "- **Data files:** %s\n", strings.Join(dataFiles, ", "))
	}
	fmt.Fprintf(&b, "- **Blog:** %s\n", yesNo(hasBlog))
	fmt.Fprintf(&b, "- **Atom feed:** %s\n", yesNo(hasFeed))
	fmt.Fprintf(&b, "- **Favicons:** %s\n", yesNo(hasFavicons))
	fmt.Fprintf(&b, "- **OG images:** %s\n", yesNo(hasOG))
	fmt.Fprintf(&b, "- **Shopify Buy button:** %s\n", yesNo(hasShopify))
	fmt.Fprintf(&b, "- **Stripe Buy button:** %s\n", yesNo(hasStripe))
	fmt.Fprintf(&b, "- **Analytics:** %s\n", yesNo(hasAnalytics))
	fmt.Fprintf(&b, "- **Chat widget:** %s\n\n", yesNo(hasChat))

	b.WriteString("## External integrations\n\n")
	b.WriteString("Opt-in snippets that call out to a third-party service. Each is wired via a `config.yaml` key and a component file whose header comment has the full setup steps — nothing else in this repo depends on these.\n\n")
	if hasShopify {
		b.WriteString("- **Shopify Buy button** — active via `site.shopify` in `config.yaml`. Setup steps: `components/shopify-buy.html`.\n")
	}
	if hasStripe {
		b.WriteString("- **Stripe Buy button** — active via `site.stripe` in `config.yaml`. Setup steps: `components/stripe-buy.html`.\n")
	}
	if hasAnalytics {
		b.WriteString("- **Analytics** — active via `site.analytics` in `config.yaml`. Setup steps: `components/analytics.html`.\n")
	}
	if hasChat {
		b.WriteString("- **Chat widget** — active via `site.chat` in `config.yaml`. Setup steps: `components/chat-widget.html`. This is an example snippet, not a bundled provider — the actual chat backend (Chatbase, Crisp, Intercom, etc.) is a separate system this repo doesn't manage.\n")
	}
	if !hasShopify && !hasStripe && !hasAnalytics && !hasChat {
		b.WriteString("- None active. Available: Shopify Buy button, Stripe Buy button, analytics (Cloudflare/Plausible/Fathom), chat widget — see `components/shopify-buy.html`, `components/stripe-buy.html`, `components/analytics.html`, `components/chat-widget.html` for setup steps.\n")
	}
	b.WriteString("\n")

	b.WriteString(`## Directory map

` + "```" + `
content/        Every file here becomes a page. content/about.md → /about.
layouts/        Page templates (Pongo2: extends/block).
components/     Reusable fragments (Pongo2: include).
data/           Global data files, exposed to templates as {{ <filename> }}.
static/         Copied to the output verbatim — CSS, images, host config.
site/           Generated output (gitignored) — this is what gets deployed.
bin/            Committed garp binaries — see Quick start above.
` + "```" + `

## Editing content

A page is Markdown with YAML frontmatter:

` + "```markdown" + `
---
title: About
layout: base.html
---

# About us
` + "```" + `

Frontmatter, a directory's ` + "`_data.yaml`" + `, and files in ` + "`data/`" + ` all merge into the variables a template sees (page frontmatter wins). Loop ` + "`collections.<section>`" + ` in a layout to list pages in a section — nothing is auto-generated.

## Commands

| Command | What it does |
|---|---|
| ` + "`bin/garp build`" + ` | Write the finished site to ` + "`site/`" + `. |
| ` + "`bin/garp serve`" + ` | Build, serve locally, rebuild on save. |
`)
	if hasFavicons {
		b.WriteString("| `bin/garp favicons <source>` | Regenerate the favicon set from a square source image. |\n")
	}
	if hasOG {
		b.WriteString("| `bin/garp og` | Regenerate the per-page OG images. |\n")
	}
	b.WriteString("| `bin/garp handoff` | Regenerate this file and refresh the committed binaries. |\n\n")

	b.WriteString("## Deploying (Cloudflare Pages)\n\n")
	b.WriteString("- **Build command:** `bin/garp build`\n")
	b.WriteString("- **Build output directory:** `site/`\n")
	b.WriteString("- **Root directory:** `/` (this repo's root)\n\n")
	b.WriteString("Push to the connected branch and Cloudflare rebuilds automatically — no environment variables or external services required; the build reads only this repo.\n")

	return b.String(), nil
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// topLevelContentDirs lists the immediate subdirectories of content/ — the
// site's "sections" (content/blog/ → "blog"), sorted.
func topLevelContentDirs(contentDir string) ([]string, error) {
	entries, err := os.ReadDir(contentDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

// topLevelFiles lists data/'s top-level yaml/yml/json files by the variable
// name they're exposed under ({{ nav }} for data/nav.yaml), sorted.
func topLevelFiles(dataDir string) ([]string, error) {
	entries, err := os.ReadDir(dataDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		switch ext {
		case ".yaml", ".yml", ".json":
			files = append(files, strings.TrimSuffix(e.Name(), ext))
		}
	}
	sort.Strings(files)
	return files, nil
}
