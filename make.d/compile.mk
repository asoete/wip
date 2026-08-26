WIP_SOURCES ?= $(shell find -type f -name "*.go" -not -path "cmd/auth-proxy-mock")
# SOURCES += Makefile
# SOURCES += $(shell find make.d/ -type f -name "*.mk")

GO := /usr/bin/go
GO_FLAGS := CGO_ENABLED=0

.PHONY: all wip
compile: wip
wip: bin/wip

# Compile a production build (excludes /dev/* and /test-fixtures/* endpoints)
bin/wip: $(WIP_SOURCES) Makefile
	$(GO_FLAGS) $(GO) build -o $@ cmd/wip/*.go
