default: build

GORELEASER := $(shell command -v goreleaser 2> /dev/null)

ifndef GORELEASER
$(error "goreleaser not found (`brew install goreleaser/tap/goreleaser` to fix)")
endif

.PHONY: build release test format

build:
	$(GORELEASER) --rm-dist --snapshot

release:
	$(GORELEASER) --rm-dist

test:
	go test -v ./...

format:
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -d {} + | tee /dev/stderr)" || \
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -w {} + | tee /dev/stderr)"
