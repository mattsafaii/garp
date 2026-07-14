package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// scaffoldFS holds the files that `garp new` stamps into a fresh project. They
// are real files under scaffold/ (not Go string consts) so they can be edited
// as normal source, byte-compared in tests, and refreshed from upstream: run
// `make sync` to re-copy the canonical safaii-css starter into
// scaffold/static/style.css. scaffold/layouts/base.html is hand-maintained
// against the safaii-html document skeleton (its <main> wrapper is load-bearing
// for the starter's content grid). The `all:` prefix keeps files that Go embed
// would otherwise skip, like a future content/blog/_data.yaml.
//
//go:embed all:scaffold
var scaffoldFS embed.FS

func cmdNew(args []string) error {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp new <path>")
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	root := fs.Arg(0)

	if _, err := os.Stat(root); err == nil {
		return fmt.Errorf("%s already exists", root)
	}

	for _, dir := range []string{"content", "layouts", "components", "data", "static", "site"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			return err
		}
	}

	if err := writeScaffold(root); err != nil {
		return err
	}

	fmt.Printf("Created %s\n\n  cd %s\n  garp dev\n", root, root)
	return nil
}

// writeScaffold copies the embedded scaffold/ tree into root, mirroring its
// structure. The scaffold/ path prefix is stripped so files land at the project
// root (scaffold/content/index.md → root/content/index.md).
func writeScaffold(root string) error {
	return fs.WalkDir(scaffoldFS, "scaffold", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel("scaffold", path)
		if err != nil {
			return err
		}
		data, err := scaffoldFS.ReadFile(path)
		if err != nil {
			return err
		}
		dst := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}
