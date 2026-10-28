# The quality gate.
#
# Every check is a target here; nobody runs `go test`, `go vet` or `gofmt` by
# hand. The gate is one command: `make check`.
#
# The tools are dev tools, not module dependencies — `make tools` installs
# them with `go install`, and `go.mod` does not move. See AGENTS.md §
# "The quality gate".

# The gate's targets touch the same files — the goldens, `go.mod`, `go.sum` —
# so a parallel `make -j` would have two of them writing at once. The gate
# takes minutes anyway; it does not need to be parallel.
.NOTPARALLEL:

GO ?= go
BIN := bin
PKG ?= ./...
RUN ?=
GOLDENS := internal/ui/testdata
STATICCHECK_VERSION := v0.8.2
GOVULNCHECK_VERSION := v1.9.0
TOOL_DOC := docs/testing.md

.PHONY: build run clean \
	fmt fmt-check vet staticcheck lint \
	test test-race \
	goldens regen-goldens goldens-check \
	budget tidy tidy-check \
	check check-race \
	tools vuln ci install-hooks versions-check

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
	GO=$(GO) scripts/check-tool-versions.sh tool staticcheck $(STATICCHECK_VERSION)
	GO=$(GO) scripts/require-tool.sh staticcheck ./...

lint:
	$(call header,$@)
	$(MAKE) fmt-check
	$(MAKE) vet
	$(MAKE) staticcheck

# One place decides how a test run is invoked. `test` and `test-race` differ
# in the flags they add and in nothing else: the filter and the package
# always reach the script, and the zero-match check always applies. Two
# recipes that spelled the invocation out separately would drift — and a slow
# gate that can pass having tested nothing is the worst kind of broken gate.
#
#   make test PKG=./internal/core RUN='-run Role|SelfReview'
#
# A filter that matches no test fails — a green that proves nothing.
TEST_FLAGS ?=
test: TEST_FLAGS :=
test-race: TEST_FLAGS := -race -shuffle=on -count=2

test test-race:
	$(call header,$@)
	GO=$(GO) scripts/go-test.sh "$(TEST_FLAGS) $(RUN)" "$(PKG)"

# The goldens are a test invocation, so they go through the same script as
# every other run: the zero-match refusal applies here too — a
# goldens-check that regenerated nothing must fail, not pass with an empty
# diff. One recipe, one spelling: the two targets differ in the diff, not
# in the run.
GOLDEN_FLAGS := -run Golden -update
GOLDEN_PKG := ./internal/ui
goldens: GOLDEN_CHECK :=
goldens-check: GOLDEN_CHECK := yes

goldens goldens-check:
	$(call header,$@)
	@if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-golden-claims.sh mark $(GOLDENS); fi
	GO=$(GO) scripts/go-test.sh "$(GOLDEN_FLAGS)" "$(GOLDEN_PKG)"
	@if [ "$(GOLDEN_CHECK)" = yes ]; then scripts/check-golden-claims.sh report $(GOLDENS); fi
	@if [ "$(GOLDEN_CHECK)" = yes ]; then \
		git diff --exit-code HEAD -- $(GOLDENS) || { \
			echo "the goldens changed — read the diff before committing it." >&2; exit 1; \
		}; \
	fi

# The old name, kept because it was in the docs and in a habit.
regen-goldens: goldens

budget:
	$(call header,$@)
	GO=$(GO) scripts/dep-budget.sh

tidy:
	$(call header,$@)
	$(GO) mod tidy

tidy-check:
	$(call header,$@)
	@tmp="$$(mktemp -d)"; cp go.mod go.sum "$$tmp/"; \
	if ! $(GO) mod tidy; then \
		cp "$$tmp/go.mod" "$$tmp/go.sum" .; \
		echo "go mod tidy failed — go.mod and go.sum were put back." >&2; \
		rm -rf "$$tmp"; exit 1; \
	fi; \
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
	$(MAKE) versions-check
	$(MAKE) test
	$(MAKE) goldens-check
	$(MAKE) budget
	$(MAKE) tidy-check

check-race:
	$(call header,$@)
	$(MAKE) lint
	$(MAKE) versions-check
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
	GO=$(GO) scripts/check-tool-versions.sh tool govulncheck $(GOVULNCHECK_VERSION)
	GO=$(GO) scripts/require-tool.sh govulncheck ./...

# The pins live here; the doc's copy is checked against them, so the doc
# cannot drift from the gate — the same shape as the dependency budget.
versions-check:
	$(call header,$@)
	scripts/check-tool-versions.sh doc $(TOOL_DOC) staticcheck $(STATICCHECK_VERSION)
	scripts/check-tool-versions.sh doc $(TOOL_DOC) govulncheck $(GOVULNCHECK_VERSION)

# What CI runs. Each check prints its own header, so the log names the one
# that failed.
ci:
	$(call header,$@)
	$(MAKE) check-race
	$(MAKE) vuln

install-hooks:
	$(call header,$@)
	@mkdir -p .git/hooks
	@if [ -e .git/hooks/pre-push ]; then \
		echo "replacing the hook that was there:"; \
		echo "---"; cat .git/hooks/pre-push; echo "---"; \
	fi
	cp scripts/hooks/pre-push.sh .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	@echo "installed .git/hooks/pre-push — it runs 'make check' before every push."
