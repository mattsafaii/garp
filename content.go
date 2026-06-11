package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Page is one Markdown file under content/.
type Page struct {
	Source string         // path relative to content/, e.g. "blog/post.md"
	Front  map[string]any // page frontmatter
	Body   string         // Markdown body with frontmatter stripped
}

func discoverContent(contentDir string) ([]*Page, error) {
	var pages []*Page
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") && path != contentDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		front, body, err := parseFrontmatter(b)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		rel, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}
		pages = append(pages, &Page{Source: rel, Front: front, Body: body})
		return nil
	})
	return pages, err
}

// parseFrontmatter splits a leading YAML frontmatter block (fenced by ---
// lines) from the body. Files without frontmatter return an empty map and
// the full source as body.
func parseFrontmatter(src []byte) (map[string]any, string, error) {
	s := string(src)
	if !strings.HasPrefix(s, "---\n") {
		return map[string]any{}, s, nil
	}
	rest := s[4:]
	var fm, body string
	if idx := strings.Index(rest, "\n---\n"); idx >= 0 {
		fm, body = rest[:idx], rest[idx+5:]
	} else if strings.HasSuffix(rest, "\n---") {
		fm, body = rest[:len(rest)-4], ""
	} else {
		return nil, "", fmt.Errorf("unclosed frontmatter")
	}

	var front map[string]any
	if err := yaml.Unmarshal([]byte(fm), &front); err != nil {
		return nil, "", fmt.Errorf("frontmatter: %w", err)
	}
	if front == nil {
		front = map[string]any{}
	}
	return front, body, nil
}
