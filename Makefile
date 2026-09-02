# -- General
# .DEFAULT_GOAL := wip.bin docs
.PHONY: default
default: wip.bin docs

###
### This file should define all variables used in multiple files
### and/or variables which can be overwritten by the user
###

# == BUILD variables ==
# ============================================================================

# -- WiP

WIP.BIN := bin/wip
WIP.SOURCES := $(shell find -type f -name "*.go" -not -path "./cmd/auth-proxy-mock/*")

# == RUN/SERVER variables ==
# ============================================================================

# -- WIP serve Config

WIP.ADDRESS := 127.0.0.1:8080
WIP.PIDFILE := /dev/shm/wip/wip.pid
DB.SQLITE_FILE := data/work-in-peace.sqlite
DB.DSN := sqlite:$(DB.SQLITE_FILE)

# -- CONTROL Config

CTL.ADDRESS := 127.0.0.1:8100
CTL.TOKEN := let-the-dev-times-roll
CTL.TOKEN_HEADER := X-WiP-Ctl-Token: $(CTL.TOKEN)

include make.d/*.mk
