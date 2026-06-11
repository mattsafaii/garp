package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const scaffoldConfig = `site_name: My Site
base_url: https://example.com
`

const scaffoldLayout = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{% if title %}{{ title }} — {% endif %}{{ site.site_name }}</title>
</head>
<body>
  {% block content %}{{ content | safe }}{% endblock %}
</body>
</html>
`

const scaffoldIndex = `---
title: Home
layout: base.html
---

# Welcome to Garp

This page is ` + "`content/index.md`" + `, rendered through ` + "`layouts/base.html`" + `.
Edit it, run ` + "`garp serve`" + `, and refresh.
`

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

	files := map[string]string{
		"config.yaml":       scaffoldConfig,
		"layouts/base.html": scaffoldLayout,
		"content/index.md":  scaffoldIndex,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			return err
		}
	}

	fmt.Printf("Created %s\n\n  cd %s\n  garp serve\n", root, root)
	return nil
}
