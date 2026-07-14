package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watchedDirs are the source dirs dev watches recursively. config.yaml is
// caught by a non-recursive watch on the project root.
var watchedDirs = []string{"content", "layouts", "components", "data", "static"}

func cmdDev(args []string) error {
	fs := flag.NewFlagSet("dev", flag.ExitOnError)
	port := fs.Int("port", 8080, "port to listen on")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: garp dev [-port N]")
	}
	fs.Parse(args)
	root := "."

	start := time.Now()
	n, err := buildSite(root)
	if err != nil {
		return err
	}
	fmt.Printf("Built %d files in %s\n", n, time.Since(start).Round(time.Millisecond))

	cfg, err := loadConfig(filepath.Join(root, "config.yaml"))
	if err != nil {
		return err
	}
	outDir := filepath.Join(root, cfg.OutputDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	if err := watcher.Add(root); err != nil {
		return err
	}
	for _, dir := range watchedDirs {
		if err := watchRecursive(watcher, filepath.Join(root, dir)); err != nil {
			return err
		}
	}
	go watchLoop(watcher, root, outDir)

	portExplicit := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			portExplicit = true
		}
	})
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil && !portExplicit {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		return err
	}
	fmt.Printf("Serving %s/ at http://localhost:%d\n", cfg.OutputDir, ln.Addr().(*net.TCPAddr).Port)
	return http.Serve(ln, siteHandler(outDir))
}

// siteHandler serves the output dir the way a static host serves it: a
// clean URL like /about falls back to about.html. Directory requests
// (/blog → blog/index.html) are http.FileServer's native behavior.
func siteHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p != "/" && path.Ext(p) == "" {
			cand := strings.TrimSuffix(p, "/") + ".html"
			if fi, err := os.Stat(filepath.Join(dir, filepath.FromSlash(cand))); err == nil && !fi.IsDir() {
				r.URL.Path = cand
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}

func watchRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == dir {
				return filepath.SkipAll
			}
			return err
		}
		if d.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

// watchLoop rebuilds on source changes, debounced — editors fire several
// events per save. Rebuild failures print and keep serving the last good
// output. Events in the output dir and hidden files are ignored; the root
// watch only counts for config.yaml.
func watchLoop(watcher *fsnotify.Watcher, root, outDir string) {
	var timer *time.Timer
	rebuild := func() {
		start := time.Now()
		n, err := buildSite(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "garp: build failed: %v\n", err)
			return
		}
		fmt.Printf("Rebuilt %d files in %s\n", n, time.Since(start).Round(time.Millisecond))
	}

	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !relevantEvent(ev.Name, root, outDir) {
				continue
			}
			if ev.Op.Has(fsnotify.Create) {
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					watchRecursive(watcher, ev.Name)
				}
			}
			if timer == nil {
				timer = time.AfterFunc(100*time.Millisecond, rebuild)
			} else {
				timer.Reset(100 * time.Millisecond)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "garp: watch error: %v\n", err)
		}
	}
}

func relevantEvent(name, root, outDir string) bool {
	if strings.HasPrefix(filepath.Base(name), ".") {
		return false
	}
	if within(outDir, name) {
		return false
	}
	// the root watch is only there for config.yaml; everything else at the
	// top level (stray files, site/ recreation) is noise
	rel, err := filepath.Rel(root, name)
	if err != nil {
		return false
	}
	if !strings.Contains(rel, string(filepath.Separator)) {
		return rel == "config.yaml" || isWatchedDir(rel)
	}
	return true
}

func isWatchedDir(rel string) bool {
	for _, d := range watchedDirs {
		if rel == d {
			return true
		}
	}
	return false
}

// within reports whether name is dir itself or inside it.
func within(dir, name string) bool {
	rel, err := filepath.Rel(dir, name)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
