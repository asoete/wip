APM.SOURCES ?= $(shell find "cmd/auth-proxy-mock" -type f -name "*.go")

APM.BIN := bin/auth-proxy-mock

APM.LISTEN_ADDR := 127.0.0.1:8888
APM.REMOTE_USER := $(USER)
APM.SSO_SUB := $(shell printf $(USER) | md5sum | awk '{ print $$1}' | sed -E 's/(.{8})(.{4})(.{4})(.{4})/\1-\2-\3-\4-/') # poor mans uuid

.PHONY: auth-proxy-mock
auth-proxy-mock: $(APM.BIN)

# Compile a production build (excludes /dev/* and /test-fixtures/* endpoints)
$(APM.BIN): $(APM.SOURCES) Makefile
	$(GO_FLAGS) $(GO) build -o $@ cmd/auth-proxy-mock/*.go

.PHONY: start-auth-proxy-mock-server
start-auth-proxy-mock-server: $(APM.BIN)
	$(APM.BIN) \
		--listen-on $(APM.LISTEN_ADDR) \
		--forward-to http://$(WEB.ADDRESS)/ \
		--remote-user $(APM.REMOTE_USER) \
		--sso-sub $(APM.SSO_SUB) \
