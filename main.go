// Command garp is a minimal static site generator.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "new":
		run(cmdNew(os.Args[2:]))
	case "build":
		run(cmdBuild(os.Args[2:]))
	case "serve":
		run(cmdServe(os.Args[2:]))
	case "favicons":
		run(cmdFavicons(os.Args[2:]))
	case "og":
		run(cmdOG(os.Args[2:]))
	case "blog":
		run(cmdBlog(os.Args[2:]))
	case "handoff":
		run(cmdHandoff(os.Args[2:]))
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "garp: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func run(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "garp: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `garp — a minimal static site generator

Usage:
  garp new <path>          Scaffold a new site
  garp build               Build the site to site/
  garp serve               Build, serve, and rebuild on change
  garp favicons <source>   Generate a favicon set from a square source image
  garp og                  Generate a templated OG image per content page
  garp blog                Scaffold an opt-in blog section + Atom feed
  garp handoff             Commit per-platform binaries + a generated HANDOFF.md
`)
}
