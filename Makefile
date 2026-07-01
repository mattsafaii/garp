# Convenience targets. Garp itself is stdlib-only; these are dev ergonomics.

STARTER := $(HOME)/.claude/skills/safaii-css/references/starter.css

.PHONY: build test sync

build:
	go build -o garp .

test:
	go test ./...

# Refresh the scaffold stylesheet from the canonical safaii-css starter, then
# verify. scaffold/layouts/base.html is hand-maintained against the safaii-html
# document skeleton and is intentionally NOT synced here.
sync:
	@test -f "$(STARTER)" || { echo "starter not found: $(STARTER)" >&2; exit 1; }
	cp "$(STARTER)" scaffold/static/style.css
	@echo "Synced scaffold/static/style.css from $(STARTER)"
	go test ./...
