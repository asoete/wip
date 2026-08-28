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

.PHONY: apm.start-server
apm.start-server: $(APM.BIN)
	$(APM.BIN) \
		--listen-on $(APM.LISTEN_ADDR) \
		--forward-to http://$(WEB.ADDRESS)/ \
		--remote-user $(APM.REMOTE_USER) \
		--sso-sub $(APM.SSO_SUB) \

.PHONY: apm.start-server.usage apm.start-server.usage.current-config apm.start-server.usage.options

apm.start-server.usage: \
	apm.start-server.usage.current-config \
	apm.start-server.usage.options

apm.start-server.usage.current-config:
	@{ \
		echo ; \
		printf "current config: \n" ; \
		printf "\t make apm.start-server APM.LISTEN_ADDR=$(APM.LISTEN_ADDR) WEB.ADDRESS=$(WEB.ADDRESS) \\ \n" ; \
		printf "\t\t APM.REMOTE_USER=$(APM.REMOTE_USER) APM.SSO_SUB=$(APM.SSO_SUB) \n" ; \
		echo ; \
	} | $(COLORIZE)

apm.start-server.usage.options:
	@{ \
		printf ' \t MAKE OPTION \t DEFAULT VALUE \t PASSED TO \t DESCRIPTION \t\n' ; \
		printf ' \t --- \t --- \t --- \t --- \t\n' ; \
		printf ' \t `APM.LISTEN_ADDR` \t `$(APM.LISTEN_ADDR)` \t `auth-proxy-mock --listen-on=...` \t Accept requests on this address \t\n' ; \
		printf ' \t `WEB.ADDRESS` \t `$(WEB.ADDRESS)` \t `auth-proxy-mock --forward-to=...` \t Forward requests to this WiP service instance \t\n' ; \
		printf ' \t `APM.REMOTE_USER` \t `$$USER` \t `auth-proxy-mock --remote-user=...` \t Inject this remote user in the HTTP request (header) \t\n' ; \
		printf ' \t `APM.SSO_SUB` \t `<md5sum($$USER) | as-uuid>` \t `auth-proxy-mock --sso-sub=...` \t Inject this sso sub in the HTTP request (header) \t\n' ; \
	} | column -t -s $$'\t' -o '|' | $(COLORIZE)
