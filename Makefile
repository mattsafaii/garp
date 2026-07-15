# Convenience targets. Garp itself is stdlib-only; these are dev ergonomics.

STARTER := $(HOME)/.claude/skills/safaii-css/references/starter.css

.PHONY: build test sync release install

build:
	go build -o garp .

# Refresh the global garp on PATH (~/go/bin) — daily-authoring convenience.
# Client repos still build with their own committed bin/garp; that pinning
# is the reproducibility guarantee, this is just Matt's local binary.
install:
	go install .

test:
	go test ./...

# Cross-compile both platforms `garp handoff` commits into a client repo:
# darwin/arm64 for Matt's own authoring, linux/amd64 for Cloudflare Pages CI.
# Run this from wherever `garp` is actually invoked from (e.g. a directory on
# PATH) so `garp handoff`, run from a client project, finds both siblings
# next to the running binary.
release:
	mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -o bin/garp-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin/garp-linux-amd64 .

# Refresh the scaffold stylesheet from the canonical safaii-css starter, then
# verify. scaffold/layouts/base.html is hand-maintained against the safaii-html
# document skeleton and is intentionally NOT synced here.
sync:
	@test -f "$(STARTER)" || { echo "starter not found: $(STARTER)" >&2; exit 1; }
	cp "$(STARTER)" scaffold/static/style.css
	@echo "Synced scaffold/static/style.css from $(STARTER)"
	go test ./...
