VERSION := $(shell git describe --abbrev=0)
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION)"

.PHONY: build
build:
	go build -v $(LDFLAGS) ./cmd/pugo

.PHONY: release
release:
	@./release.sh

.PHONY: dev-server
dev-server: build
	@cd $(dir) && ../pugo server --watch --debug --drafts

.PHONY: dev-init
dev-init: build
	@rm -rf $(dir) && mkdir -p $(dir) && cd $(dir) && ../pugo init --debug

.PHONY: dev-docs
dev-docs: build
	@cd docs && ../pugo server --watch --debug --drafts