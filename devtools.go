package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// toolbarJS is the dev toolbar, embedded so the binary stays
// self-contained. Served at /_garp/toolbar.js by dev only — it is never
// written to site/.
//
//go:embed devtools/toolbar.js
var toolbarJS []byte

func serveToolbar(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Write(toolbarJS)
}

// pageInfo is the /_garp/page response: everything the toolbar shows
// about one content page — where it came from, how it renders, and the
// resolved cascade with each key labeled by the layer that won it.
type pageInfo struct {
	URL         string      `json:"url"`
	Source      string      `json:"source"`
	LayoutChain []string    `json:"layout_chain"`
	Collections []string    `json:"collections"`
	Data        []dataEntry `json:"data"`
}

type dataEntry struct {
	Key    string `json:"key"`
	Value  any    `json:"value"`
	Source string `json:"source"`
}

// servePageInfo handles GET /_garp/page?path=<url>, computing the page's
// metadata on demand so it always reflects the current source files.
func servePageInfo(root string, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("path")
	info, err := inspectPage(root, url)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if info == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no page at " + url})
		return
	}
	json.NewEncoder(w).Encode(info)
}

// inspectPage resolves the content page served at url the same way build
// does — config, data cascade, output paths, collections — and reports
// its metadata. Returns nil (no error) when no page has that URL.
func inspectPage(root, url string) (*pageInfo, error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if err != nil {
		return nil, err
	}
	global, err := loadGlobalData(filepath.Join(root, "data"))
	if err != nil {
		return nil, err
	}
	globalSrc, err := globalDataSources(filepath.Join(root, "data"))
	if err != nil {
		return nil, err
	}
	contentDir := filepath.Join(root, "content")
	pages, err := discoverContent(contentDir)
	if err != nil {
		return nil, err
	}
	dirData, err := loadDirData(contentDir)
	if err != nil {
		return nil, err
	}

	refs := make([]*pageRef, len(pages))
	var target *pageRef
	for i, page := range pages {
		data := pageData(global, dirData, page)
		out := outputPath(page, data)
		refs[i] = &pageRef{page: page, data: data, out: out, url: pageURL(out)}
		if refs[i].url == url {
			target = refs[i]
		}
	}
	if target == nil {
		return nil, nil
	}

	layout, _ := target.data["layout"].(string)
	return &pageInfo{
		URL:         target.url,
		Source:      "content/" + filepath.ToSlash(target.page.Source),
		LayoutChain: layoutChain(root, layout),
		Collections: pageCollections(refs, target),
		Data:        cascadeEntries(cfg, globalSrc, dirData, target),
	}, nil
}

// cascadeEntries flattens a page's resolved cascade into one entry per
// top-level key, each labeled with the layer that won it (frontmatter
// beats directory data beats global). site is appended from config.yaml
// the way render exposes it; the render-time meta keys (collections,
// page, content) are never cascade data and are excluded.
func cascadeEntries(cfg *Config, globalSrc map[string]string, dirData map[string]map[string]any, ref *pageRef) []dataEntry {
	dir := filepath.Dir(ref.page.Source)
	dirSrc := path.Join("content", filepath.ToSlash(dir), "_data.yaml")

	entries := make([]dataEntry, 0, len(ref.data)+1)
	for k, v := range ref.data {
		switch k {
		case "site", "collections", "page", "content":
			continue
		}
		src := globalSrc[k]
		if _, ok := dirData[dir][k]; ok {
			src = dirSrc
		}
		if _, ok := ref.page.Front[k]; ok {
			src = "frontmatter"
		}
		entries = append(entries, dataEntry{Key: k, Value: jsonValue(v), Source: src})
	}
	entries = append(entries, dataEntry{Key: "site", Value: jsonValue(cfg.Site), Source: "config.yaml"})
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return entries
}

// globalDataSources mirrors loadGlobalData's walk of data/, mapping each
// global key to the repo-relative file it came from (data/nav.yaml → nav).
func globalDataSources(dataDir string) (map[string]string, error) {
	sources := map[string]string{}
	entries, err := os.ReadDir(dataDir)
	if os.IsNotExist(err) {
		return sources, nil
	}
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		ext := filepath.Ext(name)
		if e.IsDir() || strings.HasPrefix(name, ".") {
			continue
		}
		switch ext {
		case ".yaml", ".yml", ".json":
		default:
			continue
		}
		sources[strings.TrimSuffix(name, ext)] = "data/" + name
	}
	return sources, nil
}

// pageCollections names every collection the page belongs to, reusing
// buildCollections for membership so the answer matches what templates
// see. Sorted for a stable response.
func pageCollections(refs []*pageRef, target *pageRef) []string {
	names := []string{}
	for name, v := range buildCollections(refs) {
		entries, _ := v.([]map[string]any)
		for _, e := range entries {
			if e["url"] == target.url {
				names = append(names, name)
				break
			}
		}
	}
	sort.Strings(names)
	return names
}

var extendsRe = regexp.MustCompile(`\{%-?\s*extends\s+["']([^"']+)["']`)

// layoutChain follows a layout's extends declarations (post.html →
// base.html), resolving bare names against layouts/ then components/ the
// same way the renderer's loader does. A missing file or a cycle ends
// the chain; no layout means an empty chain.
func layoutChain(root, layout string) []string {
	chain := []string{}
	seen := map[string]bool{}
	for layout != "" && !seen[layout] {
		seen[layout] = true
		chain = append(chain, layout)
		b, err := readTemplate(root, layout)
		if err != nil {
			break
		}
		m := extendsRe.FindSubmatch(b)
		if m == nil {
			break
		}
		layout = string(m[1])
	}
	return chain
}

func readTemplate(root, name string) ([]byte, error) {
	for _, dir := range []string{"layouts", "components"} {
		if b, err := os.ReadFile(filepath.Join(root, dir, name)); err == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("template %q not found in layouts, components", name)
}

// jsonValue makes a cascade value JSON-friendly: Dates marshal as the
// ISO string templates render ("2026-05-22") instead of time.Time's
// timestamp form. Maps and slices are copied, not mutated.
func jsonValue(v any) any {
	switch t := v.(type) {
	case Date:
		return t.String()
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, val := range t {
			m[k] = jsonValue(val)
		}
		return m
	case []any:
		s := make([]any, len(t))
		for i, val := range t {
			s[i] = jsonValue(val)
		}
		return s
	}
	return v
}
