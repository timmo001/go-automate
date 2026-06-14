RM=rm -f
OUT=go-automate
TUI_OUT=go-automate-tui
DOCS_DIR=docs
VERSION=$(shell sh -c 'version=$$(git describe --long --tags --abbrev=7 2>/dev/null || printf "r%s.%s" "$$(git rev-list --count HEAD)" "$$(git rev-parse --short=7 HEAD)"); printf "%s" "$$version" | sed "s/^v//;s/\([^-]*-g\)/r\1/;s/-/./g"')

build: clean
	go build -v -ldflags="-X 'main.Version=$(VERSION)'" -o "$(OUT)" .

build_tui:
	cd tui && bun install && bun build src/index.ts --compile --outfile ../$(TUI_OUT)

build_all: build build_tui

install: create_arch
	@echo "Install with: yay -U dist/go-automate-$(VERSION)-1-x86_64.pkg.tar.zst"

run: build
	./$(OUT)

run-tui: build_tui
	./$(TUI_OUT)

docs: docs_dev

docs_dev:
	cd $(DOCS_DIR) && pnpm install && pnpm dev

docs_build:
	cd $(DOCS_DIR) && pnpm install && pnpm build

docs_preview:
	cd $(DOCS_DIR) && pnpm install && pnpm preview

docs_deps:
	cd $(DOCS_DIR) && pnpm install

test: test_go

test_go:
	go test -v ./...

lint: lint_go

lint_go:
	go fmt ./...
	go vet ./...

clean:
	-$(RM) go-automate 2>/dev/null
	-$(RM) go-automate-tui 2>/dev/null

deps:
	go mod tidy
	cd tui && bun install
	cd $(DOCS_DIR) && pnpm install

version: build
	./$(OUT) version

create_arch: clean_dist build_all
	chmod +x ./.scripts/linux/create-arch.sh
	VERSION=$(VERSION) ./.scripts/linux/create-arch.sh

clean_dist:
	-rm -rf build 2>/dev/null
	-rm -rf dist 2>/dev/null

# Show help
help:
	@echo "Available targets:"
	@echo "  build                    Build the Go application"
	@echo "  build_tui                Build the TUI binary (Bun)"
	@echo "  build_all                Build both Go app and TUI"
	@echo "  install                  Install the application"
	@echo "  run                      Build and run the application"
	@echo "  run-tui                  Build and run the TUI directly"
	@echo "  docs                     Run the documentation dev server"
	@echo "  docs_build               Build the documentation site"
	@echo "  docs_preview             Preview the built documentation site"
	@echo "  docs_deps                Install documentation dependencies"
	@echo "  test                     Run tests"
	@echo "  lint                     Run Go linters (fmt, vet)"
	@echo "  clean                    Remove build artifacts"
	@echo "  clean_dist               Remove build and dist directories"
	@echo "  deps                     Install all dependencies"
	@echo "  version                  Show the version of the application"
	@echo "  create_arch              Create Arch Linux package"
