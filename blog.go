package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// blogFS holds the files `garp blog` stamps into a project: the blog
// section (directory data, a sample post, the listing page) and the Atom
// feed pair (a frontmatter-only content page + its XML layout). Real files
// under scaffold-blog/, same convention as the scaffold/ new.go embeds.
//
//go:embed all:scaffold-blog
var blogFS embed.FS

func cmdBlog(args []string) error {
	fs := flag.NewFlagSet("blog", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp blog")
	}
	fs.Parse(args)

	written, err := stampBlog(".")
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %d files\n", len(written))
	return nil
}

// stampBlog copies the embedded scaffold-blog/ tree into root, refusing to
// write anything if any target file already exists — garp blog must be safe
// to run on an existing, populated project. Returns the paths written,
// relative to root.
func stampBlog(root string) ([]string, error) {
	var rels []string
	err := fs.WalkDir(blogFS, "scaffold-blog", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel("scaffold-blog", path)
		if err != nil {
			return err
		}
		rels = append(rels, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(rels)

	var conflicts []string
	for _, rel := range rels {
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			conflicts = append(conflicts, rel)
		}
	}
	if len(conflicts) > 0 {
		return nil, fmt.Errorf("refusing to overwrite existing file(s): %v", conflicts)
	}

	for _, rel := range rels {
		data, err := blogFS.ReadFile(filepath.Join("scaffold-blog", rel))
		if err != nil {
			return nil, err
		}
		dst := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return nil, err
		}
	}
	return rels, nil
}
