# WIP.BIN -> See main Makefile
# WIP.SOURCES -> See main Makefile

GO := /usr/bin/go
GO_FLAGS := CGO_ENABLED=0

# Compile a production build (excludes /dev/* and /test-fixtures/* endpoints)
$(WIP.BIN): $(WIP.SOURCES) Makefile
	$(GO_FLAGS) $(GO) build -o $@ cmd/wip/*.go
