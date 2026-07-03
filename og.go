package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	ogWidth  = 1200
	ogHeight = 630
	ogMargin = 80
)

func cmdOG(args []string) error {
	fs := flag.NewFlagSet("og", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp og")
	}
	fs.Parse(args)

	written, err := generateOGImages(".")
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %d OG images to static/og/\n", len(written))
	return nil
}

// generateOGImages renders one 1200x630 PNG per content page — the page
// title (wrapped) and the site name on a configurable background — and
// writes it to root/static/og/, mirroring content paths so it lines up with
// the conventional fallback path the head partial falls back to
// (/og{{ page.url }}.png). The 404 page and non-HTML permalinks (e.g. a
// feed) aren't real pages to card-share, so they're skipped, same as the
// sitemap. Returns the paths written, relative to root.
func generateOGImages(root string) ([]string, error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("config.yaml not found — is this a garp project?")
	}
	if err != nil {
		return nil, err
	}

	ogCfg, _ := cfg.Site["og"].(map[string]any)
	fontPath, _ := ogCfg["font"].(string)
	if fontPath == "" {
		return nil, fmt.Errorf(`og.font not set in config.yaml — point it at a committed .ttf, e.g.:

og:
  font: fonts/YourFont-Bold.ttf`)
	}
	fontBytes, err := os.ReadFile(filepath.Join(root, fontPath))
	if err != nil {
		return nil, fmt.Errorf("reading og.font: %w", err)
	}
	otFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing og.font %s: %w", fontPath, err)
	}

	bg := color.RGBA{R: 0x11, G: 0x11, B: 0x11, A: 0xff}
	if s, _ := ogCfg["background"].(string); s != "" {
		bg, err = parseHexColor(s)
		if err != nil {
			return nil, fmt.Errorf("og.background: %w", err)
		}
	}
	fg := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	global, err := loadGlobalData(filepath.Join(root, "data"))
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

	staticDir := filepath.Join(root, "static")
	var written []string
	for _, page := range pages {
		if page.Source == "404.md" {
			continue
		}
		data := pageData(global, dirData, page)
		out := outputPath(page, data)
		if !strings.HasSuffix(out, ".html") {
			continue
		}
		url := pageURL(out)

		title, _ := data["title"].(string)
		if title == "" {
			title = cfg.SiteName
		}
		img, err := renderOGImage(otFont, title, cfg.SiteName, bg, fg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", page.Source, err)
		}
		png, err := encodePNG(img)
		if err != nil {
			return nil, err
		}

		rel := ogRelPath(url)
		dst := filepath.Join(staticDir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dst, png, 0o644); err != nil {
			return nil, err
		}
		written = append(written, filepath.Join("static", rel))
	}
	return written, nil
}

// ogRelPath mirrors the conventional path the head partial falls back to —
// "/og{{ page.url }}.png" — as a path under static/, so the file this
// command writes always lines up with the meta tag. The home page's url is
// "/", which the raw formula would turn into the hidden file "og/.png";
// both sides special-case it to "index.png" instead.
func ogRelPath(url string) string {
	if url == "/" {
		return "og/index.png"
	}
	return "og" + url + ".png"
}

// renderOGImage draws the wrapped title (large) and site name (small) onto a
// solid-color 1200x630 canvas. Go's PNG encoder and x/image's glyph
// rasterizer are both deterministic, so the same title/font/background
// produces byte-identical output across runs.
func renderOGImage(otFont *opentype.Font, title, siteName string, bg, fg color.RGBA) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, ogWidth, ogHeight))
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	titleFace, err := opentype.NewFace(otFont, &opentype.FaceOptions{Size: 64, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, err
	}
	nameFace, err := opentype.NewFace(otFont, &opentype.FaceOptions{Size: 28, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, err
	}

	d := &font.Drawer{Dst: img, Src: image.NewUniform(fg), Face: titleFace}
	maxWidth := fixed.I(ogWidth - 2*ogMargin)
	lines := wrapText(d, title, maxWidth)

	lineHeight := titleFace.Metrics().Height
	y := fixed.I(ogMargin + 64) // baseline for the first line
	for _, line := range lines {
		d.Dot = fixed.Point26_6{X: fixed.I(ogMargin), Y: y}
		d.DrawString(line)
		y += lineHeight
	}

	d.Face = nameFace
	d.Dot = fixed.Point26_6{X: fixed.I(ogMargin), Y: fixed.I(ogHeight - ogMargin)}
	d.DrawString(siteName)

	return img, nil
}

// wrapText greedily packs words onto lines no wider than maxWidth, measured
// with the drawer's current face.
func wrapText(d *font.Drawer, text string, maxWidth fixed.Int26_6) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		trial := cur + " " + w
		if d.MeasureString(trial) > maxWidth {
			lines = append(lines, cur)
			cur = w
		} else {
			cur = trial
		}
	}
	return append(lines, cur)
}

// parseHexColor parses "#rgb" or "#rrggbb" (the # is optional).
func parseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	expand := func(c byte) string { return string([]byte{c, c}) }
	var hex string
	switch len(s) {
	case 3:
		hex = expand(s[0]) + expand(s[1]) + expand(s[2])
	case 6:
		hex = s
	default:
		return color.RGBA{}, fmt.Errorf("invalid color %q, want #rgb or #rrggbb", s)
	}
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid color %q: %w", s, err)
	}
	return color.RGBA{R: byte(v >> 16), G: byte(v >> 8), B: byte(v), A: 0xff}, nil
}
