VERSION    := $(shell git describe --tags --always 2>/dev/null || echo "0.0.0-dev")
COMMIT     := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -X 'gostim2/internal/version.Version=$(VERSION)' \
              -X 'gostim2/internal/version.GitCommit=$(COMMIT)' \
              -X 'gostim2/internal/version.BuildTime=$(BUILD_TIME)'

PLATFORMS  := linux/amd64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64
DIST       := dist

.DEFAULT_GOAL := help

.PHONY: help build build-multiplatform clean test fmt vet run run-gui install

## help: Show this help message
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'

## build: Build CLI and GUI binaries for the current platform
build:
	go build -ldflags "$(LDFLAGS)" -o gostim2 ./cmd/gostim2
	go build -ldflags "$(LDFLAGS)" -o gostim2-gui ./cmd/gostim2-gui
	@echo "Built gostim2 and gostim2-gui ($(VERSION))"

## build-multiplatform: Build CLI and GUI for all target platforms into dist/
build-multiplatform:
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d/ -f1); \
		arch=$$(echo $$platform | cut -d/ -f2); \
		ext=; [ "$$os" = "windows" ] && ext=.exe; \
		dir=$(DIST)/gostim2-$(VERSION)-$$os-$$arch; \
		echo "Building $$os/$$arch..."; \
		mkdir -p $$dir; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $$dir/gostim2$$ext ./cmd/gostim2 || exit 1; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $$dir/gostim2-gui$$ext ./cmd/gostim2-gui || exit 1; \
	done
	@echo "All builds written to $(DIST)/"

## install: Install binaries to GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/gostim2
	go install -ldflags "$(LDFLAGS)" ./cmd/gostim2-gui

## install-completion: Install bash and zsh completions for gostim2
install-completion:
	./scripts/install-completions.sh

## test: Run all tests
test:
	go test ./...

## fmt: Format all Go source files
fmt:
	go fmt ./...

## vet: Run go vet on all packages
vet:
	go vet ./...

## run: Run the CLI (pass ARGS="..." to supply arguments)
run:
	go run -ldflags "$(LDFLAGS)" ./cmd/gostim2 $(ARGS)

## run-gui: Run the GUI (pass ARGS="..." to supply arguments)
run-gui:
	go run -ldflags "$(LDFLAGS)" ./cmd/gostim2-gui $(ARGS)

## clean: Remove built binaries and dist/
clean:
	rm -f gostim2 gostim2-gui
	rm -rf $(DIST)
