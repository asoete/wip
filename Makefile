SOURCES ?= $(shell find -type f -name "*.go")

GO := /usr/bin/go
GO_FLAGS := CGO_ENABLED=0

PORT := 8080
CTL_PORT := 8100
CTL_TOKEN := let-the-dev-times-roll 123

.PHONY: all
all: bin/wip

# Compile a production build (excludes /dev/* and /test-fixtures/* endpoints)
bin/wip: $(SOURCES) Makefile
	$(GO_FLAGS) $(GO) build -o $@ cmd/wip/*.go

# Compile and run a DEV build
.PHONY: run
run: fmt
	WIP_CTL_TOKEN="$(CTL_TOKEN)" \
		$(GO_FLAGS) $(GO) run -tags debug cmd/wip/*.go -port $(PORT)

.PHONY: fmt
fmt:
	@echo $(dir $(SOURCES)) | tr ' ' '\n' | sort | uniq | xargs -I {} go fmt {}

# Kill a running server via the builtin controls
.PHONY: kill-server
kill-server:
	@printf ">_ MAKE: !! KILLING SERVER !! [POST http://localhost:$(CTL_PORT)/ctl/exit/0]\n"
	@curl -s -X POST -H "X-WiP-Ctl-Token: $(CTL_TOKEN)" localhost:$(CTL_PORT)/ctl/exit/0

include make.d/*.mk
