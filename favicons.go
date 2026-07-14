package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

// icoSizes are the frames packed into favicon.ico, the conventional trio
// browsers pick from for the tab/bookmark icon.
var icoSizes = []int{16, 32, 48}

func cmdFavicons(args []string) error {
	fs := flag.NewFlagSet("favicons", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp favicons <source>")
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}

	written, err := generateFavicons(".", fs.Arg(0))
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %d files to static/\n", len(written))
	return nil
}

// generateFavicons reads a single square source image and writes the full
// favicon set into root/static/. Pure-Go resize (golang.org/x/image/draw,
// CatmullRom) and Go's deterministic PNG encoder mean the same source
// produces byte-identical output across runs — clean diffs when a client's
// source logo hasn't changed. Returns the paths written, relative to root.
func generateFavicons(root, source string) ([]string, error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("config.yaml not found — is this a garp project?")
	}
	if err != nil {
		return nil, err
	}

	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("reading source image: %w", err)
	}
	src, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		return nil, fmt.Errorf("decoding %s: %w", source, err)
	}

	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w != h {
		return nil, fmt.Errorf("source image must be square, got %dx%d", w, h)
	}
	if w < 1024 {
		return nil, fmt.Errorf("source image must be at least 1024px, got %dx%d", w, h)
	}

	staticDir := filepath.Join(root, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		return nil, err
	}

	var written []string
	write := func(name string, data []byte) error {
		if err := os.WriteFile(filepath.Join(staticDir, name), data, 0o644); err != nil {
			return err
		}
		written = append(written, filepath.Join("static", name))
		return nil
	}

	icoFrames := make([]image.Image, len(icoSizes))
	for i, size := range icoSizes {
		icoFrames[i] = resizeSquare(src, size)
	}
	ico, err := encodeICO(icoFrames)
	if err != nil {
		return nil, err
	}
	if err := write("favicon.ico", ico); err != nil {
		return nil, err
	}

	for name, size := range map[string]int{"icon-192.png": 192, "icon-512.png": 512, "apple-touch-icon.png": 180} {
		png, err := encodePNG(resizeSquare(src, size))
		if err != nil {
			return nil, err
		}
		if err := write(name, png); err != nil {
			return nil, err
		}
	}

	manifest, err := encodeManifest(cfg.SiteName)
	if err != nil {
		return nil, err
	}
	if err := write("site.webmanifest", manifest); err != nil {
		return nil, err
	}

	return written, nil
}

// resizeSquare scales a square source image to size x size using CatmullRom,
// a high-quality interpolating filter well suited to shrinking a large logo
// down to icon sizes.
func resizeSquare(src image.Image, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return dst
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encodeICO packs images into a .ico container using the PNG-in-ICO format
// (supported since Windows Vista, and by every browser favicon.ico consumer
// worth targeting) — no separate ICO-codec dependency needed. Format:
// https://en.wikipedia.org/wiki/ICO_(file_format)
func encodeICO(images []image.Image) ([]byte, error) {
	frames := make([][]byte, len(images))
	for i, img := range images {
		b, err := encodePNG(img)
		if err != nil {
			return nil, err
		}
		frames[i] = b
	}

	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&out, binary.LittleEndian, uint16(1)) // type: icon
	binary.Write(&out, binary.LittleEndian, uint16(len(images)))

	offset := uint32(6 + 16*len(images))
	for i, img := range images {
		b := img.Bounds()
		// a byte holds 0-255; ICO's convention is that a 256px side wraps
		// to 0, which byte(256) does for free
		out.WriteByte(byte(b.Dx()))
		out.WriteByte(byte(b.Dy()))
		out.WriteByte(0)                                    // color count (0 = not palette-indexed)
		out.WriteByte(0)                                    // reserved
		binary.Write(&out, binary.LittleEndian, uint16(1))  // planes
		binary.Write(&out, binary.LittleEndian, uint16(32)) // bits per pixel
		binary.Write(&out, binary.LittleEndian, uint32(len(frames[i])))
		binary.Write(&out, binary.LittleEndian, offset)
		offset += uint32(len(frames[i]))
	}
	for _, frame := range frames {
		out.Write(frame)
	}
	return out.Bytes(), nil
}

type webManifest struct {
	Name            string         `json:"name"`
	ShortName       string         `json:"short_name"`
	Icons           []manifestIcon `json:"icons"`
	ThemeColor      string         `json:"theme_color"`
	BackgroundColor string         `json:"background_color"`
	Display         string         `json:"display"`
	StartURL        string         `json:"start_url"`
}

type manifestIcon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

func encodeManifest(siteName string) ([]byte, error) {
	m := webManifest{
		Name:      siteName,
		ShortName: siteName,
		Icons: []manifestIcon{
			{Src: "/icon-192.png", Sizes: "192x192", Type: "image/png"},
			{Src: "/icon-512.png", Sizes: "512x512", Type: "image/png"},
		},
		ThemeColor:      "#ffffff",
		BackgroundColor: "#ffffff",
		Display:         "standalone",
		// Without start_url, installs open at whatever page linked the manifest.
		StartURL: "/",
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}
