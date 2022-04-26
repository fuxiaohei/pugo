VERSION := $(shell cat ./version)
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION)"


.PHONY: build
build:
	go build -v $(LDFLAGS) ./cmd/pugo

.PHONY: release
release:
	@./release.sh