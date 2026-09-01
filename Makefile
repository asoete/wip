# -- General
# .DEFAULT_GOAL := wip.bin docs
.PHONY: default
default: wip.bin docs

# -- MAIN Config

WEB.ADDRESS := 127.0.0.1:8080
WIP.PIDFILE := /dev/shm/wip/wip.pid
DB.SQLITE_FILE := data/work-in-peace.sqlite
DB.DSN := sqlite:$(DB.SQLITE_FILE)

# -- CONTROL Config

CTL.ADDRESS := 127.0.0.1:8100
CTL.TOKEN := let-the-dev-times-roll
CTL.TOKEN_HEADER := X-WiP-Ctl-Token: $(CTL.TOKEN)

# -- Building
WIP.SOURCES ?= $(shell find -type f -name "*.go" -not -path "./cmd/auth-proxy-mock/*")

include make.d/*.mk
