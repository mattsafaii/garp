package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// loadGlobalData reads the top-level files of data/ — each file becomes one
// template variable named after the file: data/nav.yaml → nav. A missing
// data/ dir is fine.
func loadGlobalData(dataDir string) (map[string]any, error) {
	global := map[string]any{}
	entries, err := os.ReadDir(dataDir)
	if os.IsNotExist(err) {
		return global, nil
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
		b, err := os.ReadFile(filepath.Join(dataDir, name))
		if err != nil {
			return nil, err
		}
		var v any
		if err := yaml.Unmarshal(b, &v); err != nil {
			return nil, fmt.Errorf("data/%s: %w", name, err)
		}
		global[strings.TrimSuffix(name, ext)] = v
	}
	return global, nil
}

// loadDirData finds _data.yaml files under content/ and returns them keyed by
// directory relative to content/ ("." for the root). Each applies to pages in
// that directory only.
func loadDirData(contentDir string) (map[string]map[string]any, error) {
	dirData := map[string]map[string]any{}
	err := filepath.WalkDir(contentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "_data.yaml" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var v map[string]any
		if err := yaml.Unmarshal(b, &v); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		rel, err := filepath.Rel(contentDir, filepath.Dir(path))
		if err != nil {
			return err
		}
		dirData[rel] = v
		return nil
	})
	if os.IsNotExist(err) {
		return dirData, nil
	}
	return dirData, err
}

// pageData merges the cascade for one page: global → its directory's
// _data.yaml → frontmatter, shallow, last wins.
func pageData(global map[string]any, dirData map[string]map[string]any, page *Page) map[string]any {
	merged := map[string]any{}
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range dirData[filepath.Dir(page.Source)] {
		merged[k] = v
	}
	for k, v := range page.Front {
		merged[k] = v
	}
	return merged
}
