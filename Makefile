RM=rm -f
OUT=go-automate
VERSION=$(shell sh -c 'version=$$(git describe --long --tags --abbrev=7 2>/dev/null || printf "r%s.%s" "$$(git rev-list --count HEAD)" "$$(git rev-parse --short=7 HEAD)"); printf "%s" "$$version" | sed "s/^v//;s/\([^-]*-g\)/r\1/;s/-/./g"')

build: clean
	go build -v -ldflags="-X 'main.Version=$(VERSION)'" -o "$(OUT)" .

install: create_arch
	@echo "Install with: yay -U dist/go-automate-$(VERSION)-1-x86_64.pkg.tar.zst"

run: build
	./$(OUT)

test: test_go

test_go:
	go test -v ./...

lint: lint_go

lint_go:
	go fmt ./...
	go vet ./...

clean:
	-$(RM) go-automate 2>/dev/null

deps:
	go mod tidy

version: build
	./$(OUT) version

create_arch: clean_dist build
	chmod +x ./.scripts/linux/create-arch.sh
	VERSION=$(VERSION) ./.scripts/linux/create-arch.sh

clean_dist:
	-rm -rf build 2>/dev/null
	-rm -rf dist 2>/dev/null

# Show help
help:
	@echo "Available targets:"
	@echo "  build                    Build the application"
	@echo "  install                  Install the application"
	@echo "  run                      Build and run the application"
	@echo "  test                     Run tests"
	@echo "  lint                     Run Go linters (fmt, vet)"
	@echo "  clean                    Remove build artifacts"
	@echo "  clean_dist               Remove build and dist directories"
	@echo "  deps                     Install dependencies"
	@echo "  version                  Show the version of the application"
	@echo "  create_arch              Create Arch Linux package"
