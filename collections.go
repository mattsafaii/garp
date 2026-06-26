package main

import (
	"sort"
	"strings"
	"time"
)

// pageRef is a page with its cascade resolved — what build computes before
// rendering and what collections are built from.
type pageRef struct {
	page *Page
	data map[string]any
	out  string // output path under site/
	url  string // clean URL the host serves
}

// buildCollections gathers pages into collections.* variables. Every page is
// in collections.all. A page joins collections.<dir> for its top-level
// content directory (content/blog/post.md → collections.blog) — except the
// directory's own index.md, which is the listing page, not a member. tags:
// in the cascade (string or list) adds the page to collections.<tag>.
// Entries sort newest-first by date; undated entries follow, by path.
func buildCollections(refs []*pageRef) map[string]any {
	groups := map[string][]*pageRef{}
	seen := map[string]map[*pageRef]bool{}
	add := func(name string, r *pageRef) {
		if seen[name] == nil {
			seen[name] = map[*pageRef]bool{}
		}
		if seen[name][r] {
			return
		}
		seen[name][r] = true
		groups[name] = append(groups[name], r)
	}

	for _, r := range refs {
		add("all", r)
		if dir, rest, ok := strings.Cut(r.page.Source, "/"); ok && rest != "index.md" {
			add(dir, r)
		}
		for _, tag := range tagList(r.data) {
			add(tag, r)
		}
	}

	collections := map[string]any{}
	for name, refs := range groups {
		sort.SliceStable(refs, func(i, j int) bool {
			di, dj := pageDate(refs[i].data), pageDate(refs[j].data)
			if !di.Equal(dj) {
				return di.After(dj)
			}
			return refs[i].page.Source < refs[j].page.Source
		})
		entries := make([]map[string]any, len(refs))
		for i, r := range refs {
			entries[i] = map[string]any{
				"url":   r.url,
				"data":  r.data,
				"title": r.data["title"],
				"date":  r.data["date"],
			}
		}
		collections[name] = entries
	}
	return collections
}

func tagList(data map[string]any) []string {
	switch v := data["tags"].(type) {
	case string:
		return []string{v}
	case []any:
		var tags []string
		for _, t := range v {
			if s, ok := t.(string); ok {
				tags = append(tags, s)
			}
		}
		return tags
	}
	return nil
}

func pageDate(data map[string]any) time.Time {
	switch t := data["date"].(type) {
	case Date:
		return t.Time
	case time.Time:
		return t
	}
	return time.Time{}
}
