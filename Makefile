SOURCES ?= $(shell find -type f -name "*.go")

GO := /usr/bin/go
GO_FLAGS := CGO_ENABLED=0

bin/wip: $(SOURCES) Makefile
	$(GO_FLAGS) time -p $(GO) build -o $@ cmd/wip/main.go

run:
	$(GO_FLAGS) $(GO) run cmd/wip/main.go

