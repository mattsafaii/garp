# Garp

Garp is a fast, minimal static site generator written in Go. It compiles to a single dependency-free binary and exists to fill one gap: **Nunjucks-style templating in a single binary, with no Node toolchain.** Hugo is a single binary but its Go templating is miserable; 11ty has the templating and data model but needs Node/npm. Garp gives you Pongo2 (Jinja2/Nunjucks-familiar) templates and an 11ty-like data cascade and collections, emitting plain static files that deploy to any host. It's built to craft client marketing sites for Safaii Studio and Zonebrite.

## Stack

- **Language:** Go (single binary, no runtime dependencies for the end user). Compiles for Mac/Linux only — Windows may compile but is not tested or supported.
- **Dependencies:** Pongo2 (templating), goldmark (Markdown), fsnotify (file watching), gopkg.in/yaml.v3 (config + frontmatter). Don't add others without a clear reason.
- **CLI:** **stdlib `flag` only — do NOT use Cobra or any CLI framework.** Three commands dispatch via `switch os.Args[1]`, each with its own `flag.FlagSet` for flags. Revisit this only if commands ever exceed ~8 or gain nested subcommands (the no-gos say they won't).

## Commands (surfaces)

- **`garp new <path>`** — scaffolds a project: the six reserved dirs (`content/ layouts/ components/ data/ static/ site/`), a `config.yaml`, one base layout, and a sample `content/index.md` that renders immediately with zero edits.
- **`garp build`** — reads `config.yaml`, walks `content/`, merges the data cascade, renders Markdown + Pongo2 with layout chaining, writes flat `.html` to `site/`, copies `static/` verbatim. Prints file count + build time.
- **`garp serve`** — runs build, serves `site/` over local HTTP, watches `content/ layouts/ components/ data/ static/ config.yaml` via fsnotify, rebuilds on change. Prints the local URL. Takes `-port N` (default 8080); without an explicit `-port` it falls back to an OS-assigned port if 8080 is taken. The HTTP handler mirrors the host: a clean URL like `/about` falls back to `about.html`.

## Conventions (hard rules)

- **Reserved directories:** `content/ layouts/ components/ data/ static/ site/`. Only `content/` produces output pages; everything else is machinery and is never output. `site/` is the generated output (gitignored).
- **URL/permalink mapping:** 1:1 mirror, flat `.html`, **no** index-in-a-directory expansion. `content/about.md` → `site/about.html` → served at `/about`; `content/blog/post.md` → `/blog/post`; `content/index.md` → `/`. A `permalink:` in frontmatter overrides the output path. (Clean no-trailing-slash URLs are the host's job — Cloudflare Pages serves `about.html` at `/about`. The generator does nothing host-specific.)
- **Data cascade:** three levels only — global (`data/`) → directory data → page frontmatter, last wins. No computed data, no JS data files, no data-from-functions.
  - **Global data** is namespaced by filename: `data/nav.yaml` → `{{ nav }}`, `data/company.json` → `{{ company.* }}`. Only top-level `.yaml`/`.yml`/`.json` files in `data/` are loaded (not nested dirs).
  - **Directory data** is a `_data.yaml` inside a content directory, applying to pages in **that** directory only (not subdirectories). `content/blog/_data.yaml` is the idiomatic place to set `layout: post.html` for a whole section.
  - The merge is shallow (top-level keys, last wins). Frontmatter beats directory data beats global.
- **Config:** single `config.yaml`. Known keys `site_name`, `base_url`, `output_dir`; any additional key the author adds is exposed to every template under `site.*`. Minimal but extensible.
- **Collections are variables, not pages.** Gather content into `collections.*` exposed to templates. The generator NEVER auto-generates archive/taxonomy pages. A blog index is just `content/blog/index.md` with a layout that loops over `collections.blog`. The available collections:
  - `collections.all` — every page.
  - **One per top-level content directory** — `content/blog/post.md` joins `collections.blog`. The directory's own `index.md` is excluded (it's the listing page, not a member).
  - **One per `tags:` entry** (string or list in the cascade) — a page with `tags: [featured]` joins `collections.featured`.
  - Entries carry `url`, `title`, `date`, and full `data`. Sorted newest-first by `date`; undated entries follow, ordered by source path. A page reached two ways (its directory plus a tag) is deduped within each collection.
- **Templates:** Pongo2, using `extends`/`block`/`include`. Layouts live in `layouts/`, reusable fragments in `components/`. Both `extends "x.html"` and `include "x.html"` resolve by bare name against `layouts/` then `components/` (no path prefixes). A page's rendered HTML reaches its layout as `{{ content }}` — emit it with `{{ content | safe }}`. Page Markdown bodies are run through Pongo2 *before* goldmark, so `{{ site.* }}` and cascade variables work inside content too.

## No-gos (do not build these)

- Swift Mac companion app / CMS (deferred to v2/v3 — not a line of Swift this cycle).
- Auto-generated taxonomy / archive pages.
- Built-in pagination.
- Host-specific deploy adapters (no Cloudflare/Netlify/GitHub Pages modes). Host config files like `_redirects`, `_headers`, `.nojekyll` are just authored in `static/` and copied verbatim.
- Asset pipeline — no Sass, no JS bundling, no minification. Vanilla CSS copied as a static file.
- Image processing / responsive images.
- Plugin or extensibility API.
- Incremental / cached builds — full rebuild every time.
- Multiple template engines — Pongo2 only.
- Cobra / CLI frameworks — stdlib `flag` only.
- Windows support.

## Philosophy

Do the obvious thing. No magic, no unnecessary abstractions. Single binary. Convention over configuration. Zero client-side JS shipped by default. Build speed is a first-class feature.

## Basecamp

This project's Basecamp config (account / project / todolist IDs) is already set in `.basecamp/config.json` (gitignored), so `basecamp` commands work without flags from this directory. The work is the **Build** todolist — run `basecamp todos list` to see it, and check each todo off as you complete it. The full PRD is a Basecamp doc: run `basecamp docs list` and open the one titled **PRD**. The shaped pitch lives on the project's card in the Lab Ideas board.

## Solo

This project runs inside Solo, which manages its processes. For coordination and follow-up work, call `help(topic="coordination")` and `help(topic="timers")`.
