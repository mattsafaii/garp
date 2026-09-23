# Garp

A fast, minimal static site generator in a single Go binary.

Garp fills one gap: **Nunjucks-style templating in a single binary, with no Node toolchain.** Hugo is a single binary but its Go templating is painful to write; 11ty has the templating and data model you want but needs Node/npm. Garp pairs [Pongo2](https://github.com/flosch/pongo2) (Django/Jinja2/Nunjucks-familiar) templates with an 11ty-like data cascade and collections, and emits plain static HTML that deploys to any host.

It's a personal tool, built for crafting Safaii Studio client marketing sites.

## Install

```sh
go install github.com/mattsafaii/garp@latest
```

Requires Go 1.26+. Builds for macOS and Linux. (Windows may compile but is untested and unsupported.)

## Quick start

```sh
garp new mysite     # scaffold a project that builds with zero edits
cd mysite
garp dev            # build, serve at http://localhost:8080, rebuild + reload on save
```

Edit a file, save, refresh. When you're ready to ship:

```sh
garp build          # write the finished site to site/
```

Point your host at `site/` and you're live.

## Commands

| Command | What it does |
|---|---|
| `garp new <path>` | Scaffold a project: the six reserved dirs, a `config.yaml`, a base layout, and a sample `content/index.md` that renders immediately. |
| `garp build` | Read `config.yaml`, walk `content/`, merge the data cascade, render Markdown + Pongo2 with layout chaining, write flat `.html` to `site/`, copy `static/` verbatim. With `sink: true` in config, also synthesizes `kitchen-sink.html` — a noindex design-system reference page generated from the stylesheet's `@layer tokens` (author-overridable, like sitemap.xml). Prints the file count and build time. |
| `garp dev [-port N]` | Run a build, serve `site/` over local HTTP, watch the source via fsnotify, rebuild on change, and live-reload the browser. Includes a dev toolbar (the "g" button, bottom-right) that inspects the current page's data cascade — each key labeled with the layer that set it — plus its layout chain and collections, and a kitchen-sink page at `/_garp/sink` regenerated from the stylesheet's tokens on every request. Injected into dev responses only, never into `site/`. Defaults to port 8080, falling back to an OS-assigned port if it's taken. |
| `garp favicons <source>` | Generate a full favicon set (`favicon.ico`, `icon-192.png`, `icon-512.png`, `apple-touch-icon.png`, `site.webmanifest`) from one square source image, into `static/`. |
| `garp og` | Generate a templated 1200×630 OG image per content page, into `static/og/`. |
| `garp blog` | Stamp an opt-in blog section — a post layout, a listing page, a sample post, and an Atom feed. Refuses to overwrite existing files. |
| `garp handoff` | Write committed per-platform binaries (`bin/`) plus a generated `HANDOFF.md` — makes the repo buildable and deployable by anyone, without garp installed. |

## How a project is laid out

```
mysite/
├── config.yaml     # site_name, base_url, output_dir + any custom keys
├── content/        # the ONLY directory that becomes pages
├── layouts/        # page templates (extends / block)
├── components/     # reusable fragments (include)
├── data/           # global data files
├── static/         # copied to the output verbatim (incl. the seed style.css)
└── site/           # generated output (gitignored)
```

These six directory names are reserved. Only `content/` produces output pages; everything else is machinery.

## Writing pages

A page is Markdown with optional YAML frontmatter. `layout:` names a file in `layouts/`:

```markdown
---
title: About
layout: base.html
---

# About us

We do {{ site.site_name }} things.
```

Paths mirror 1:1, flat: `content/about.md` → `site/about.html` (served at `/about`), `content/blog/post.md` → `/blog/post`, `content/index.md` → `/`. A `permalink:` in frontmatter overrides the output path.

Markdown bodies are run through Pongo2 *before* Markdown conversion, so `{{ site.* }}` and data values work inside content too.

## The data cascade

Every page renders with a merged bag of variables, built from three layers — last wins:

1. **Global** — each file in `data/` becomes a variable named after it: `data/nav.yaml` → `{{ nav }}`.
2. **Directory** — a `_data.yaml` in a content directory applies to pages in that directory (handy for setting `layout:` once for a whole section).
3. **Page frontmatter** — wins over both.

No computed data, no JavaScript data files — three predictable layers.

## Collections

Garp never auto-generates archive or tag pages. Instead it gathers pages into `collections.*` variables you loop over yourself:

- `collections.all` — every page.
- `collections.<dir>` — every page in a top-level content directory (the directory's own `index.md` is excluded — it's the listing page).
- `collections.<tag>` — every page with that `tags:` entry.

Entries carry `url`, `title`, `date`, and the full page `data`, sorted newest-first by date. A blog index is just a normal page:

```html
{% for post in collections.blog %}
  <a href="{{ post.url }}">{{ post.title }}</a> — {{ post.date | date:"January 2, 2006" }}
{% endfor %}
```

## Templates

Pongo2, with `extends` / `block` / `include`. Layouts live in `layouts/`, fragments in `components/`; both resolve by bare name (`{% include "footer.html" %}`). A page's rendered HTML reaches its layout as `{{ content }}` — emit it with `{{ content | safe }}`.

## Dates

A YAML date (`date: 2026-05-22`) renders as a clean ISO date by default (`{{ date }}` → `2026-05-22`). Format it with the `date` filter and a Go reference layout:

```
{{ post.date | date:"January 2, 2006" }}   →  May 22, 2026
```

## Config

A single `config.yaml`. The known keys are `site_name`, `base_url`, and `output_dir` (defaults to `site`). Any additional key you add is exposed to every template under `site.*`.

```yaml
site_name: Acme
base_url: https://acme.com
phone: "(555) 555-0148"   # → {{ site.phone }}
```

## Deploying

Garp does nothing host-specific. It writes plain static files; clean no-trailing-slash URLs are the host's job (e.g. Cloudflare Pages serves `about.html` at `/about`). Host config files like `_headers` or `_redirects` go in `static/` and are copied through untouched.

## Toolbelt

Opt-in commands for recurring chores, run at authoring time (never in CI) with their artifacts committed.

**`garp favicons <source>`** — reads one square source image (1024px+) and writes `favicon.ico`, `icon-192.png`, `icon-512.png`, `apple-touch-icon.png`, and `site.webmanifest` into `static/`. The scaffold's head partial already links them; rerun any time the source logo changes — output is byte-stable for the same input.

**`garp og`** — renders a 1200×630 PNG per page (title + site name on a solid background) into `static/og/`, mirroring each page's URL. Needs a font pointed at from `config.yaml`:

```yaml
og:
  font: fonts/YourFont-Bold.ttf   # a committed .ttf — nothing is embedded
  background: "#111111"           # optional, defaults to a dark gray
```

With `og:` set, any page without an `image:` in frontmatter gets the matching OG/Twitter meta tags automatically; an explicit `image:` still overrides.

**`garp blog`** — stamps an opt-in blog section into a new or existing project: `content/blog/_data.yaml`, `layouts/post.html`, a sample post, a listing page, and an Atom feed (`content/feed.md` + `layouts/feed.xml`, served at `/feed.xml`). Refuses to run if any target file already exists.

**`garp handoff`** — writes `bin/` (the running binary plus any other `garp-<goos>-<goarch>` binaries sitting next to it — build both with `make release` in the garp repo) and a generated `HANDOFF.md` describing the project's actual shape. Safe to rerun any time. Point Cloudflare Pages at build command `bin/garp build`, output directory `site/`.

## Built with

[Pongo2](https://github.com/flosch/pongo2) (templating), [goldmark](https://github.com/yuin/goldmark) (Markdown), [fsnotify](https://github.com/fsnotify/fsnotify) (file watching), [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) (config + frontmatter + data), and [golang.org/x/image](https://pkg.go.dev/golang.org/x/image) (favicon resizing + OG text rendering, toolbelt-only). The CLI is the Go standard library only.
