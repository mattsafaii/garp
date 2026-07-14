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
	"sync"
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
	rl := newReloader()
	go watchLoop(watcher, root, outDir, rl)

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
	return http.Serve(ln, siteHandler(root, outDir, rl))
}

// toolbarTag rides along on every HTML response dev serves — it is never
// written to site/. The script itself (live reload + the page inspector
// toolbar) is embedded in the binary and served at /_garp/toolbar.js.
const toolbarTag = `<script src="/_garp/toolbar.js" defer></script>` + "\n"

// reloader broadcasts rebuild events to connected browsers over SSE;
// the injected reloadScript listens and reloads the page.
type reloader struct {
	mu   sync.Mutex
	subs map[chan struct{}]struct{}
}

func newReloader() *reloader {
	return &reloader{subs: make(map[chan struct{}]struct{})}
}

func (rl *reloader) broadcast() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for ch := range rl.subs {
		select {
		case ch <- struct{}{}:
		default: // subscriber already has a pending event
		}
	}
}

func (rl *reloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	ch := make(chan struct{}, 1)
	rl.mu.Lock()
	rl.subs[ch] = struct{}{}
	rl.mu.Unlock()
	defer func() {
		rl.mu.Lock()
		delete(rl.subs, ch)
		rl.mu.Unlock()
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprint(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}

// siteHandler serves the output dir the way a static host serves it: a
// clean URL like /about falls back to about.html, and a path that
// resolves to nothing gets the project's 404.html with a 404 status —
// the same way Cloudflare Pages serves it in production. Directory
// requests (/blog → blog/index.html) are http.FileServer's native
// behavior. HTML responses get the dev toolbar script tag appended;
// everything else streams through FileServer untouched. root is the
// project root — dev-only endpoints under /_garp/ read source files
// (content/, data/, layouts/) the output dir no longer carries.
func siteHandler(root, dir string, rl *reloader) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/_garp/") {
			switch r.URL.Path {
			case "/_garp/reload":
				rl.ServeHTTP(w, r)
			case "/_garp/toolbar.js":
				serveToolbar(w)
			case "/_garp/page":
				servePageInfo(root, w, r)
			default:
				http.NotFound(w, r)
			}
			return
		}
		p := r.URL.Path
		if p != "/" && path.Ext(p) == "" {
			cand := strings.TrimSuffix(p, "/") + ".html"
			if fi, err := os.Stat(filepath.Join(dir, filepath.FromSlash(cand))); err == nil && !fi.IsDir() {
				r.URL.Path = cand
			}
		}
		if name, ok := htmlTarget(dir, r.URL.Path); ok {
			serveHTML(w, http.StatusOK, name)
			return
		}
		if onDisk(dir, r.URL.Path) {
			fileServer.ServeHTTP(w, r)
			return
		}
		if serveHTML(w, http.StatusNotFound, filepath.Join(dir, "404.html")) {
			return
		}
		// no 404 page in this project; FileServer's plain 404 will do
		fileServer.ServeHTTP(w, r)
	})
}

// htmlTarget resolves a request path to the .html file dev serves by
// hand: direct .html requests, and directory requests' index.html ("/"
// or a trailing slash). Bare directory paths (/blog) stay with
// FileServer for its canonical trailing-slash redirect.
func htmlTarget(dir, p string) (string, bool) {
	fsp := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+p)))
	switch {
	case strings.HasSuffix(p, ".html"):
	case p == "/" || strings.HasSuffix(p, "/"):
		fsp = filepath.Join(fsp, "index.html")
	default:
		return "", false
	}
	if fi, err := os.Stat(fsp); err == nil && !fi.IsDir() {
		return fsp, true
	}
	return "", false
}

// serveHTML writes an HTML file with the dev toolbar script tag appended.
func serveHTML(w http.ResponseWriter, status int, name string) bool {
	body, err := os.ReadFile(name)
	if err != nil {
		return false
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
	w.Write([]byte(toolbarTag))
	return true
}

// onDisk reports whether a request path resolves to an existing file or
// directory under dir (directories are FileServer's job: redirects,
// index.html).
func onDisk(dir, p string) bool {
	_, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+p))))
	return err == nil
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
func watchLoop(watcher *fsnotify.Watcher, root, outDir string, rl *reloader) {
	var timer *time.Timer
	rebuild := func() {
		start := time.Now()
		n, err := buildSite(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "garp: build failed: %v\n", err)
			return
		}
		fmt.Printf("Rebuilt %d files in %s\n", n, time.Since(start).Round(time.Millisecond))
		rl.broadcast()
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
