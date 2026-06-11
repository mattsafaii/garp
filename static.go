package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// copyStatic copies static/ into the output dir verbatim, preserving the
// directory structure. Returns the number of files copied. A missing
// static/ dir is fine.
func copyStatic(staticDir, outDir string) (int, error) {
	count := 0
	err := filepath.WalkDir(staticDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(outDir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := copyFile(path, dst); err != nil {
			return err
		}
		count++
		return nil
	})
	if os.IsNotExist(err) {
		return 0, nil
	}
	return count, err
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
