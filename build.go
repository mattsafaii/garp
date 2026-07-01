package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func cmdBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp build")
	}
	fs.Parse(args)

	start := time.Now()
	n, err := buildSite(".")
	if err != nil {
		return err
	}
	fmt.Printf("Built %d files in %s\n", n, time.Since(start).Round(time.Millisecond))
	return nil
}

// buildSite runs the full pipeline for the project at root: config, data
// cascade, content, collections, render, static passthrough. The output dir
// is wiped first — it's generated, disposable output. Returns total files
// written.
func buildSite(root string) (int, error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return 0, fmt.Errorf("config.yaml not found — is this a garp project?")
	}
	if err != nil {
		return 0, err
	}

	global, err := loadGlobalData(filepath.Join(root, "data"))
	if err != nil {
		return 0, err
	}
	contentDir := filepath.Join(root, "content")
	pages, err := discoverContent(contentDir)
	if err != nil {
		return 0, err
	}
	dirData, err := loadDirData(contentDir)
	if err != nil {
		return 0, err
	}

	refs := make([]*pageRef, len(pages))
	for i, page := range pages {
		data := pageData(global, dirData, page)
		out := outputPath(page, data)
		refs[i] = &pageRef{page: page, data: data, out: out, url: pageURL(out)}
	}
	collections := buildCollections(refs)

	outDir := filepath.Join(root, cfg.OutputDir)
	if err := checkOutputDir(root, outDir, cfg.OutputDir); err != nil {
		return 0, err
	}
	if err := os.RemoveAll(outDir); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return 0, err
	}

	r := newRenderer(filepath.Join(root, "layouts"), filepath.Join(root, "components"))
	for _, ref := range refs {
		// fresh map — collections entries hold ref.data, which must not
		// grow a reference back to collections
		ctx := make(map[string]any, len(ref.data)+3)
		for k, v := range ref.data {
			ctx[k] = v
		}
		ctx["site"] = cfg.Site
		ctx["collections"] = collections
		ctx["page"] = map[string]any{"url": ref.url, "source": ref.page.Source}
		html, err := r.renderPage(ref.page, ctx)
		if err != nil {
			return 0, err
		}
		dst := filepath.Join(outDir, ref.out)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return 0, err
		}
		if err := os.WriteFile(dst, []byte(html), 0o644); err != nil {
			return 0, err
		}
	}

	staticCount, err := copyStatic(filepath.Join(root, "static"), outDir)
	if err != nil {
		return 0, err
	}

	seoCount, err := synthesizeSEO(root, outDir, cfg, refs)
	if err != nil {
		return 0, err
	}

	return len(refs) + staticCount + seoCount, nil
}

// checkOutputDir rejects an output_dir the build would wipe destructively:
// the project root itself, anything outside it, or a source directory —
// os.RemoveAll(outDir) must only ever hit generated output.
func checkOutputDir(root, outDir, configured string) error {
	rel, err := filepath.Rel(root, outDir)
	if err != nil || rel == "." {
		return fmt.Errorf("output_dir %q would overwrite the project itself", configured)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("output_dir %q is outside the project", configured)
	}
	top, _, _ := strings.Cut(rel, string(filepath.Separator))
	for _, d := range []string{"content", "layouts", "components", "data", "static"} {
		if top == d {
			return fmt.Errorf("output_dir %q would delete the %s/ source directory", configured, d)
		}
	}
	return nil
}

// synthesizeSEO writes sitemap.xml and robots.txt into outDir from the page
// refs and base_url. Either is author-overridable: if static/sitemap.xml or
// static/robots.txt already exists, copyStatic has just placed the author's
// own version in outDir and synthesis for that file is skipped so it isn't
// clobbered.
func synthesizeSEO(root, outDir string, cfg *Config, refs []*pageRef) (int, error) {
	n := 0
	if _, err := os.Stat(filepath.Join(root, "static", "sitemap.xml")); os.IsNotExist(err) {
		if err := os.WriteFile(filepath.Join(outDir, "sitemap.xml"), []byte(renderSitemap(cfg.BaseURL, refs)), 0o644); err != nil {
			return n, err
		}
		n++
	} else if err != nil {
		return n, err
	}

	if _, err := os.Stat(filepath.Join(root, "static", "robots.txt")); os.IsNotExist(err) {
		if err := os.WriteFile(filepath.Join(outDir, "robots.txt"), []byte(renderRobots(cfg.BaseURL)), 0o644); err != nil {
			return n, err
		}
		n++
	} else if err != nil {
		return n, err
	}
	return n, nil
}

// renderSitemap lists every page ref as an absolute URL under base_url, with
// lastmod from the page's frontmatter date when it has one. The 404 page is
// excluded — it isn't a real destination to index.
func renderSitemap(baseURL string, refs []*pageRef) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, ref := range refs {
		if ref.page.Source == "404.md" {
			continue
		}
		b.WriteString("\t<url>\n")
		b.WriteString("\t\t<loc>" + xmlEscape(baseURL+ref.url) + "</loc>\n")
		if d := pageDate(ref.data); !d.IsZero() {
			b.WriteString("\t\t<lastmod>" + d.Format("2006-01-02") + "</lastmod>\n")
		}
		b.WriteString("\t</url>\n")
	}
	b.WriteString("</urlset>\n")
	return b.String()
}

func renderRobots(baseURL string) string {
	return "User-agent: *\nAllow: /\n\nSitemap: " + baseURL + "/sitemap.xml\n"
}

func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}
