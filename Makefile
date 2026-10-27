GO ?= go
BIN := bin

.PHONY: build test lint fmt clean run regen-goldens

build:
	@mkdir -p $(BIN)
	$(GO) build -o $(BIN)/scale ./cmd/scale
	$(GO) build -o $(BIN)/maitre ./cmd/maitre

test:
	$(GO) test ./...

lint:
	@test -z "$$(gofmt -l . )" || (echo "gofmt needed:"; gofmt -l .; exit 1)
	$(GO) vet ./...

fmt:
	gofmt -w .

run:
	$(GO) run ./cmd/scale

regen-goldens:
	$(GO) test ./internal/ui -run Golden -update

clean:
	rm -rf $(BIN)
	$(GO) clean -testcache
