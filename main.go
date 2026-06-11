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
  garp new <path>   Scaffold a new site
  garp build        Build the site to site/
  garp serve        Build, serve, and rebuild on change
`)
}
