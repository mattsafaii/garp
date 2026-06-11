package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	if outDir == filepath.Clean(root) {
		return 0, fmt.Errorf("output_dir %q would overwrite the project itself", cfg.OutputDir)
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
	return len(refs) + staticCount, nil
}
