package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
)

// dirsLoader resolves template names against a list of directories, first
// hit wins. Names stay logical ("base.html") rather than being made
// absolute, so extends/include in layouts/ can pull from components/ and
// vice versa without path prefixes.
type dirsLoader struct {
	dirs []string
}

func (l *dirsLoader) Abs(base, name string) string { return name }

func (l *dirsLoader) Get(path string) (io.Reader, error) {
	for _, dir := range l.dirs {
		if b, err := os.ReadFile(filepath.Join(dir, path)); err == nil {
			return bytes.NewReader(b), nil
		}
	}
	return nil, fmt.Errorf("template %q not found in %s", path, strings.Join(l.dirs, ", "))
}

// renderer wraps a Pongo2 template set over layouts/ and components/ — so
// {% extends "base.html" %} and {% include "footer.html" %} both work
// without path prefixes.
type renderer struct {
	set *pongo2.TemplateSet
}

func newRenderer(layoutsDir, componentsDir string) *renderer {
	return &renderer{set: pongo2.NewSet("garp", &dirsLoader{dirs: []string{layoutsDir, componentsDir}})}
}

// renderPage runs a page's Markdown body through Pongo2 (so {{ site.* }} and
// friends resolve inside content), converts it to HTML, then wraps it in the
// layout named by the cascade — the layout receives the HTML as content.
// Pages without a layout emit the bare Markdown HTML.
func (r *renderer) renderPage(page *Page, ctx map[string]any) (string, error) {
	bodyTpl, err := r.set.FromString(page.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", page.Source, err)
	}
	rendered, err := bodyTpl.Execute(pongo2.Context(ctx))
	if err != nil {
		return "", fmt.Errorf("%s: %w", page.Source, err)
	}
	html, err := renderMarkdown(rendered)
	if err != nil {
		return "", fmt.Errorf("%s: %w", page.Source, err)
	}

	layout, _ := ctx["layout"].(string)
	if layout == "" {
		return html, nil
	}
	tpl, err := r.set.FromCache(layout)
	if err != nil {
		return "", fmt.Errorf("%s: layout %q: %w", page.Source, layout, err)
	}
	layoutCtx := pongo2.Context{}
	for k, v := range ctx {
		layoutCtx[k] = v
	}
	layoutCtx["content"] = html
	out, err := tpl.Execute(layoutCtx)
	if err != nil {
		return "", fmt.Errorf("%s: layout %q: %w", page.Source, layout, err)
	}
	return out, nil
}
