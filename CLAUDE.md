# Garp

Garp is a fast, minimal static site generator written in Go. It compiles to a single dependency-free binary and exists to fill one gap: **Nunjucks-style templating in a single binary, with no Node toolchain.** Hugo is a single binary but its Go templating is miserable; 11ty has the templating and data model but needs Node/npm. Garp gives you Pongo2 (Jinja2/Nunjucks-familiar) templates and an 11ty-like data cascade and collections, emitting plain static files that deploy to any host.

Garp is a **private studio tool** — Matt's kit for building and maintaining Safaii Studio client sites. It is never distributed as a product: no public repo, no docs site, no community, no support surface. Its purpose is to make Matt better at the actual website work instead of fighting tooling.

## What garp is (and isn't)

**The compiled output is the product; the tooling vanishes from what ships.** A website is just HTML + CSS + JS. Garp's job is to generate that plain artifact and then get out of the way. Every design decision serves that: the deliverable must be maintainable by any developer with zero garp knowledge, must outlive garp itself, and must never carry studio tooling as baggage.

Garp is organized in three layers. Keep them separate — the separation is what keeps the durable core small.

- **Engine (Go)** — `build` / `dev` / `new`. Tiny, fast, dependency-free, afternoon-readable. This is the deliverable-maker and the thing whose longevity everything depends on. It changes rarely and stays legible on purpose (bus-factor). Already built.
- **Toolbelt (Go)** — opt-in commands that automate recurring client chores (SEO files, favicons, OG images, handoff bundles). Each is subject to the litmus test below. Grows phase by phase through the roadmap below — nothing speculative beyond it. Toolbelt commands live in their own files/packages and must never bloat the engine.
- **Cockpit (Swift, later)** — a separate native Mac app ("my own Framer") that *drives* the garp binary; it never replaces it. Native SwiftUI + iCloud for Matt's own multi-site management and editing. Deferred; its own project, not bolted onto an engine cycle. The engine stays **Go** because the build must run portably in CI and on any future developer's machine — a Swift build tool would lock buildability to macOS and destroy the handoff/bus-factor guarantee.

## Constraints

1. **Litmus test.** Every tool garp gains must leave behind a plain artifact that survives without garp. If a feature would make the *output* depend on garp to function, it doesn't belong.
2. **Handoff = repo + committed binary.** Clients receive the full source repo plus the compiled garp binary (optionally vendored garp source for the "outlive me forever" case). **Source is always committed; output-only handoff is never done.** This is deliberately *anti*-lock-in: a client can leave, or Matt can go MIA, and any developer picks the site up. A single dependency-free Go binary makes this handoff *more* durable than an npm-based SSG repo, not less — it builds with one command, no toolchain, forever.
3. **CI-purity line.** Anything pure-Go (sitemap, favicons, OG text, SEO files) may run inside the Cloudflare Pages build via the committed binary. Anything needing an external tool (search indexing → pagefind, image optimization → vips/cwebp, or cgo) runs on **Matt's Mac at authoring time**, with artifacts committed. The CI build must never hard-depend on an external binary being installed, or the deploy becomes fragile and the bus-factor guarantee breaks.
4. **Sync is git.** The site lives in git; git is the versioning, multi-person, and deploy-trigger substrate (Cloudflare Pages builds from it). Never put site content under iCloud/Drive sync — a git working tree inside iCloud corrupts. iCloud is only ever for the future cockpit's own app state.
5. **Build step ≠ toolchain.** The build step is inherent and correct — it's the one thing the web platform never solved (build-time templating/DRY). Garp keeps the *step* and eliminated the *toolchain* (no node_modules, no version managers). Run in CI via the committed binary, the build step is invisible: edit content, push, the site rebuilds. Never treat the build step as baggage to remove.

## Stack

- **Language:** Go (single binary, no runtime dependencies for the end user). Compiles for Mac/Linux only — Windows may compile but is not tested or supported.
- **Engine dependencies:** Pongo2 (templating), goldmark (Markdown), fsnotify (file watching — `dev` only), gopkg.in/yaml.v3 (config + frontmatter). Don't add others to the engine without a clear reason. Vendor dependencies (`go mod vendor`) so garp builds offline and survives upstream deletion — combined with Go's compat promise, a vendored garp compiles for decades.
- **Toolbelt dependencies:** a toolbelt command may pull an additional *Go* library (compiled into the binary — categorically different from npm baggage). External binaries (pagefind, vips/cwebp) are allowed only for Mac-authoring-time commands, never for CI-path commands. Still ask before adding.
- **CLI:** **stdlib `flag` only — do NOT use Cobra or any CLI framework.** Commands dispatch via `switch os.Args[1]`, each with its own `flag.FlagSet`. The flat command list will grow as the toolbelt fills in; that's fine — a bigger switch is not a reason for Cobra. Reconsider only if commands ever gain genuinely *nested* subcommands (they shouldn't).

## Commands

### Engine (built, stable)

- **`garp new <path>`** — scaffolds a project: the six reserved dirs (`content/ layouts/ components/ data/ static/ site/`), a `config.yaml`, one base layout, and a sample `content/index.md` that renders immediately with zero edits. This scaffold is also where most SEO/head/security defaults ship (see roadmap).
- **`garp build`** — reads `config.yaml`, walks `content/`, merges the data cascade, renders Markdown + Pongo2 with layout chaining, writes flat `.html` to `site/`, copies `static/` verbatim. Also synthesizes `sitemap.xml` and `robots.txt` from the page refs and `base_url`; both are author-overridable — a file already present in `static/` always wins and synthesis for it is skipped. Prints file count + build time.
- **`garp dev`** — runs build, serves `site/` over local HTTP, watches `content/ layouts/ components/ data/ static/ config.yaml` via fsnotify, rebuilds on change, and live-reloads the browser: an SSE endpoint at `/_garp/reload` plus a one-line script appended to HTML *responses* only — `site/` on disk is never touched, so the litmus test isn't grazed. Prints the local URL. Takes `-port N` (default 8080); without an explicit `-port` it falls back to an OS-assigned port if 8080 is taken. The HTTP handler mirrors the host: a clean URL like `/about` falls back to `about.html`, and missing paths get the project's `404.html` with a 404 status, like Cloudflare Pages. (Renamed from `garp serve` 2026-07-13 — local development only, never a production server.)

### Toolbelt (Phase 2, shipped 2026-07-02)

Each is a separate command in its own file, subject to the litmus test. See the roadmap for phase order and sequencing.

- **`garp favicons <source>`** — reads one square source image (1024px+; errors if not square or too small) and emits `favicon.ico` (a hand-rolled PNG-in-ICO container, 16/32/48px frames), `icon-192.png`, `icon-512.png`, `apple-touch-icon.png` (180×180), and `site.webmanifest` into `static/`. Pure-Go resize via `golang.org/x/image/draw` (CatmullRom); Go's PNG encoder is deterministic, so output is byte-stable across runs on the same source. The scaffold's `components/favicons.html` (included from `head.html`) carries the corresponding `<link>` tags unconditionally — harmless 404s until the command has been run once.
- **`garp og`** — renders a templated 1200×630 PNG per content page (wrapped title + site name on a solid, configurable background) and writes it to `static/og/<path>.png`, mirroring the page's URL (the home page's is `static/og/index.png`). Text comes from a committed `.ttf` named in config (`og.font`) — errors with a clear message when unset; nothing is embedded in the binary. `components/head.html` falls back to the conventional path `/og{{ page.url }}.png` when `site.og` is set and a page has no explicit `image:` frontmatter (which always wins). Deterministic PNG + glyph rendering means unchanged pages produce identical bytes.
- **`garp blog`** — stamps an opt-in blog section: `content/blog/_data.yaml` (`layout: post.html`), `layouts/post.html`, a sample rendering post (`content/blog/hello-world.md`), a listing `content/blog/index.md` looping `collections.blog`, and the Atom feed pair — a frontmatter-only `content/feed.md` (`permalink: /feed.xml`, `layout: feed.xml`, empty body) plus `layouts/feed.xml`, a pure Pongo2 Atom template that ignores `{{ content }}` and loops `collections.blog` directly (the empty body means goldmark never touches the XML). Works on a new or existing project; refuses to write anything if any target file already exists (checked up front, all-or-nothing).
- **`garp handoff`** — the bus-factor command. Writes `bin/` — the running binary plus any `garp-<goos>-<goarch>` siblings found next to it (so `bin/garp-darwin-arm64` + `bin/garp-linux-amd64`, built via `make release` in the garp repo, both land when handoff is run from a Mac binary that has the Linux one sitting alongside it) — plus a three-line `bin/garp` shell shim that dispatches on `uname`, and a generated `HANDOFF.md` describing *this* project's actual shape (sections, data files, which toolbelt commands are in use) with a command reference and Cloudflare Pages settings (build command `bin/garp build`, output dir `site/`). Safe to rerun any time — it never touches `content/`/`layouts/`/`data/`/`static/`, only `bin/` and `HANDOFF.md`. Warns (non-fatal) rather than failing when the Linux sibling isn't found next to the running binary.

### Post-Phase-2 additions (shipped 2026-07-07)

Added ad hoc from a dogfooding conversation, not a formal phase — passes the litmus test on its own terms.

- **Speculation Rules prefetch** (added 2026-07-13, from the Maudit comparison) — `components/prefetch.html`, included from `head.html` and gated like analytics: emits nothing unless `prefetch: true` is set in `config.yaml`, then emits a `<script type="speculationrules">` block prefetching same-site links on hover/intent (`eagerness: moderate`). A platform primitive — no JS ships; unsupported browsers ignore the script type.
- **`garp search`** — the first command needing an external binary outside Phase 3's planned image-optimization fork: shells out to the `pagefind` CLI (https://pagefind.app) against the already-built site, writing the index straight to `static/pagefind/` (via pagefind's own `--output-path` flag) instead of the gitignored `site/pagefind/`. Mac-authoring-time only, per the CI-purity line — errors clearly if `site/` hasn't been built yet or `pagefind` isn't on `PATH`. Because the index lands in `static/`, a plain `garp build` (no pagefind installed) reproduces it in CI on every deploy — same authoring-time-then-commit pattern as favicons/OG. The scaffold's opt-in `components/search-ui.html` wires up pagefind's own prebuilt `PagefindUI` widget; nothing to build ourselves.

## Conventions (hard rules)

- **Reserved directories:** `content/ layouts/ components/ data/ static/ site/`. Only `content/` produces output pages; everything else is machinery and is never output. `site/` is the generated output (gitignored).
- **URL/permalink mapping:** 1:1 mirror, flat `.html`, **no** index-in-a-directory expansion. `content/about.md` → `site/about.html` → served at `/about`; `content/blog/post.md` → `/blog/post`; `content/index.md` → `/`. A `permalink:` in frontmatter overrides the output path. (Clean no-trailing-slash URLs are the host's job — Cloudflare Pages serves `about.html` at `/about`. The generator does nothing host-specific.)
- **Data cascade:** three levels only — global (`data/`) → directory data → page frontmatter, last wins. No computed data, no JS data files, no data-from-functions. `garp build` reads only static local data files — this keeps builds deterministic, offline, and reproducible from the repo alone (a bus-factor guarantee). **Remote sources (Airtable, a third-party API, a SQLite DB) are supported out of band:** an authoring-time `garp fetch` command snapshots them into committed static data files, which the cascade then consumes like any other data. Fetching never happens inside `build` or in CI — so a client's site always rebuilds from the repo alone, with no API keys or live services. Fetched data populates template **variables**; it may never drive a generated page set (see no-gos).
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
- **Dates:** a YAML date in frontmatter or data (`date: 2026-05-22`) is wrapped so a bare `{{ date }}` renders as an ISO date (`2026-05-22`), not Go's raw `time.Time` string. Format it with the `date` filter and a Go reference layout: `{{ post.date | date:"January 2, 2006" }}` → `May 22, 2026`; `{{ post.date | date }}` with no argument falls back to ISO. (`time` is an alias for the same filter.) Collections still sort by the underlying date value.

## Roadmap (build incrementally)

Phases, not deadlines. Sequencing (decided 2026-07-01, Phase 2 shipped 2026-07-02): the toolbelt is now complete — next up is dogfooding the whole kit on the first real client site before any Phase 3 work starts. Stick to the commands listed; nothing speculative beyond the roadmap.

- **Phase 1 (shipped 2026-07-01)** — unblocks the small-catalog Shopify client and every future site. No new commands; a richer `garp new` scaffold plus one small build-time step. Scaffold cluster: `<head>` meta/OG/Twitter partial (with canonical, via the existing `page.url`), 404 page, `_headers` security defaults, config-driven analytics partial (Cloudflare Web Analytics / Plausible / Fathom), JSON-LD LocalBusiness partial + `data/business.yaml`, client-side Shopify Buy SDK snippet. Build step: synthesize **sitemap.xml** and **robots.txt** (from the page refs + `base_url`; both author-overridable — if the file exists in `static/`, the author's wins).
- **Phase 2 (shipped 2026-07-02 — see the "PRD: Phase 2 — Toolbelt" Basecamp doc for the full decisions)** — a **favicons** command (pure-Go resize; adds `golang.org/x/image` for quality downscaling), templated **OG images** (pure-Go text-on-image; font is a config-pointed committed TTF, nothing embedded), an opt-in **`garp blog`** scaffold command, an **Atom feed** (authored template stamped by `garp blog`, layout-driven — never build-synthesized), **`garp handoff`** (makes the repo self-sufficient: per-platform committed binaries + a generated per-project HANDOFF.md — the bus-factor command; no zip), a **privacy-policy page** in the scaffold (no cookie banner — the supported analytics are cookieless), and a **Stripe buy button** snippet. See the Toolbelt section above for what each command actually does.
- **Phase 3 / later** — **`garp fetch`** (snapshot remote data from Airtable / a third-party API / a SQLite DB into committed static data files → the cascade reads them; Mac-authoring-time, never in CI — see the data cascade rule), **image optimization** (the one genuine dependency fork: modern-format encoding needs vips/cwebp or cgo → Mac-authoring-time, artifacts committed), an **`llms.txt` synthesis** step alongside the existing sitemap.xml/robots.txt build-step synthesis (pure-Go, deterministic, author-overridable the same way), a scaffolded **design-system reference page** rendering the project's CSS tokens (colors, type scale, spacing) as a plain static page for client/dev reference, the client self-edit **CMS** track (git-based web CMS like PagesCMS, wired by a command — not built into the binary), and the **Swift cockpit**.

## No-gos

### Permanent — never in the engine

- Auto-generated taxonomy / archive pages. A page set generated from a data collection — local *or* fetched — is out. Fetched data may populate template variables; it may never emit one-page-per-item. This is the line that keeps large-catalog ecommerce out of garp.
- Built-in pagination.
- Remote / computed / function-derived data **inside the build** — `build` and the cascade read only static local files. Remote data enters only via the out-of-band `garp fetch` snapshot command (see roadmap), which writes static files the cascade then reads. What stays forbidden is `build` itself reaching the network or running data functions.
- Asset pipeline in the build — no Sass, no JS bundling, no minification. Vanilla CSS copied as a static file. (Image optimization is a *Mac-authoring-time toolbelt command*, never a CI build step — see phase 3.)
- Plugin or extensibility API.
- Incremental / cached builds — full rebuild every time.
- Multiple template engines — Pongo2 only.
- Cobra / CLI frameworks — stdlib `flag` only.
- Windows support.

### Reclassified (were flat no-gos; now planned, out of the durable core)

- **Swift Mac app / CMS** — no longer forbidden, but it is the **cockpit layer** (phase 3+), a separate app that drives the binary. Still not part of an engine cycle.
- **Image processing / responsive images** — a phase 3 **toolbelt** command that runs at Mac-authoring-time and commits plain artifacts. Never in the engine, never in the CI build path.

### The governing line

Not in the durable core; allowed as an opt-in toolbelt (or the cockpit), and only if it passes the litmus test — it must leave a plain artifact that survives without garp. Host config files (`_redirects`, `_headers`, `.nojekyll`) and third-party snippets (Shopify Buy SDK, analytics) are just authored/scaffolded into `static/` or templates and shipped verbatim — that's allowed and is not a "host-specific deploy adapter." Ecommerce default is **client-side Shopify**; large catalogs needing generated product pages go to Shopify native, not garp.

## Philosophy

Do the obvious thing. No magic, no unnecessary abstractions. Single binary. Convention over configuration. Zero client-side JS shipped by default (opt-in snippets like Shopify or analytics are the author's explicit choice). Build speed is a first-class feature.

**AI-accelerated, never AI-dependent.** Agents and skills may drive garp; garp, its output, and its handoffs never reference or require AI. The dependency arrow points one way only — and no convention may live solely in an AI layer: human-readable docs (README, HANDOFF.md) are the source of truth, a skill is at most a translation of them. If AI disappeared tomorrow, everything still works; the only thing lost is speed.

Above all: **the tooling vanishes into a self-sufficient deliverable, and the engine stays small enough that any developer — including a future you — can read the whole thing in an afternoon.**

## Basecamp

This project's Basecamp config (account / project / todolist IDs) is already set in `.basecamp/config.json` (gitignored), so `basecamp` commands work without flags from this directory.

**Work is always tracked in Basecamp** — it's where the project, todos, and progress live. Run `basecamp todos list` to see the active todolist and check items off as you complete them. Each build cycle/phase gets its own todolist; the original **Build** list (v1), **Phase 1 — SEO scaffold + Shopify snippet**, **Post-review fixes — 2026-07-02**, and **Phase 2 — Toolbelt** are all complete. There is no active cycle right now — per the roadmap sequencing, next up is dogfooding the whole kit on the first real Safaii Studio client site, not a new phase. The Phase 2 spec is the doc titled **PRD: Phase 2 — Toolbelt** (run `basecamp docs list`), including the two design notes decided at scoping (layout-driven feed, per-platform committed binaries). The v1 PRD is the doc titled **PRD**; the shaped pitch lives on the project's card in the Lab Ideas board. Architecture/decisions/module-level documentation for the whole codebase lives in the project's **Documentation** folder in Docs & Files (kept current via `/document`).

## Solo

This project runs inside Solo, which manages its processes. For coordination and follow-up work, call `help(topic="coordination")` and `help(topic="timers")`.
