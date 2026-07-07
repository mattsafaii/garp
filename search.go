package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func cmdSearch(args []string) error {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp search")
	}
	fs.Parse(args)

	n, err := runSearchIndex(".")
	if err != nil {
		return err
	}
	fmt.Printf("Indexed %d files into static/pagefind/\n", n)
	return nil
}

// runSearchIndex shells out to the pagefind CLI (https://pagefind.app) — an
// external Rust binary, so per the CI-purity line this only ever runs at
// Mac-authoring-time, never in CI. It writes the index straight into
// static/pagefind/ (via pagefind's own --output-path flag) rather than
// site/pagefind/, so a plain `garp build` reproduces the search bundle on
// every future build/deploy without pagefind installed — the same
// authoring-time-then-commit pattern as favicons and OG images. Returns the
// number of files pagefind wrote.
func runSearchIndex(root string) (int, error) {
	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return 0, fmt.Errorf("config.yaml not found — is this a garp project?")
	}
	if err != nil {
		return 0, err
	}

	if _, err := os.Stat(filepath.Join(root, cfg.OutputDir)); errors.Is(err, fs.ErrNotExist) {
		return 0, fmt.Errorf("%s not found — run `garp build` first", cfg.OutputDir)
	}

	if _, err := exec.LookPath("pagefind"); err != nil {
		return 0, fmt.Errorf(`pagefind not found on PATH — install it (https://pagefind.app/docs/installation/), e.g.:

npm install -g pagefind

then rerun garp search`)
	}

	outputPath := filepath.Join("static", "pagefind")
	cmd := exec.Command("pagefind", "--site", cfg.OutputDir, "--output-path", outputPath)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("running pagefind: %w", err)
	}

	n := 0
	err = filepath.WalkDir(filepath.Join(root, outputPath), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			n++
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("pagefind ran but %s wasn't created: %w", outputPath, err)
	}
	return n, nil
}
