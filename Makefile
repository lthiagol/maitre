# The quality gate.
#
# Every check is a target here; nobody runs `go test`, `go vet` or `gofmt` by
# hand. The gate is one command: `make check`.
#
# The tools are dev tools, not module dependencies — `make tools` installs
# them with `go install`, and `go.mod` does not move. See AGENTS.md §
# "The quality gate".

GO ?= go
BIN := bin
PKG ?= ./...
RUN ?=
GOLDENS := internal/ui/testdata
STATICCHECK_VERSION := v0.8.2
GOVULNCHECK_VERSION := v1.9.0

.PHONY: build run clean \
	fmt fmt-check vet staticcheck lint \
	test test-race \
	goldens regen-goldens goldens-check \
	budget tidy tidy-check \
	check check-race \
	tools vuln ci install-hooks

define header
@echo ""
@echo "==> $(1)"
endef

build:
	$(call header,$@)
	@mkdir -p $(BIN)
	$(GO) build -o $(BIN)/maitre-cli ./cmd/maitre-cli
	$(GO) build -o $(BIN)/maitre ./cmd/maitre

run:
	$(call header,$@)
	$(GO) run ./cmd/maitre-cli

clean:
	$(call header,$@)
	rm -rf $(BIN)
	$(GO) clean -testcache

# --- the checks -----------------------------------------------------------

fmt:
	$(call header,$@)
	gofmt -w .

fmt-check:
	$(call header,$@)
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then \
		echo "gofmt needed:"; echo "$$out"; exit 1; \
	fi

vet:
	$(call header,$@)
	$(GO) vet ./...

staticcheck:
	$(call header,$@)
	scripts/require-tool.sh staticcheck ./...

lint:
	$(call header,$@)
	$(MAKE) fmt-check
	$(MAKE) vet
	$(MAKE) staticcheck

# `PKG` and `RUN` narrow the run:
#   make test PKG=./internal/core RUN='-run Role|SelfReview'
# A filter that matches no test fails — a green that proves nothing.
test:
	$(call header,$@)
	GO=$(GO) scripts/go-test.sh "$(RUN)" "$(PKG)"

test-race:
	$(call header,$@)
	GO=$(GO) scripts/go-test.sh "-race -shuffle=on -count=2" "$(PKG)"

goldens:
	$(call header,$@)
	$(GO) test ./internal/ui -run Golden -update

# The old name, kept because it was in the docs and in a habit.
regen-goldens: goldens

goldens-check:
	$(call header,$@)
	$(GO) test ./internal/ui -run Golden -update
	@git diff --exit-code HEAD -- $(GOLDENS) || { \
		echo "the goldens changed — read the diff before committing it." >&2; exit 1; \
	}

budget:
	$(call header,$@)
	scripts/dep-budget.sh

tidy:
	$(call header,$@)
	$(GO) mod tidy

tidy-check:
	$(call header,$@)
	@tmp="$$(mktemp -d)"; cp go.mod go.sum "$$tmp/"; \
	$(GO) mod tidy; \
	if ! git diff --quiet HEAD -- go.mod go.sum; then \
		cp "$$tmp/go.mod" "$$tmp/go.sum" .; \
		echo "go.mod or go.sum is not tidy — run 'make tidy' and commit the result." >&2; \
		rm -rf "$$tmp"; exit 1; \
	fi; \
	rm -rf "$$tmp"

# --- the gates ------------------------------------------------------------

check:
	$(call header,$@)
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) goldens-check
	$(MAKE) budget
	$(MAKE) tidy-check

check-race:
	$(call header,$@)
	$(MAKE) lint
	$(MAKE) test-race
	$(MAKE) goldens-check
	$(MAKE) budget
	$(MAKE) tidy-check

# --- the tools ------------------------------------------------------------

tools:
	$(call header,$@)
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	$(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	@echo ""
	@echo "staticcheck $(STATICCHECK_VERSION) and govulncheck $(GOVULNCHECK_VERSION) are installed."
	@echo "binaries: $$(go env GOPATH)/bin"
	@case ":$$PATH:" in \
		*":$$(go env GOPATH)/bin:"*) ;; \
		*) echo "That directory is not on your PATH — add it, or the tools stay out of reach." ;; \
	esac

vuln:
	$(call header,$@)
	scripts/require-tool.sh govulncheck ./...

# What CI runs. Each check prints its own header, so the log names the one
# that failed.
ci:
	$(call header,$@)
	$(MAKE) check-race
	$(MAKE) vuln

install-hooks:
	$(call header,$@)
	@mkdir -p .git/hooks
	cp scripts/hooks/pre-push.sh .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	@echo "installed .git/hooks/pre-push — it runs 'make check' before every push."
