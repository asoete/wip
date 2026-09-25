APM.SOURCES ?= $(shell find "cmd/auth-proxy-mock" -type f -name "*.go")

APM.BIN := bin/auth-proxy-mock

APM.LISTEN_ADDR := 127.0.0.1:8888
APM.USERNAME := $(USER)
APM.REMOTE_USER := $(shell printf $(USER) | ./tools/stdin-to-uuid.sh)@sso.irc.ugent.be/auth/realms/irc # poor mans uuid
APM.GROUPS := users,it,services,vib,wwwauth,wwwu3,wwwu8,bioit,wwwcm,plasman,wwwgbe2,wwwubla,bioitadmins,mybit,u_ysa,collaboratest,s_bioit,sc-core,itdevel

.PHONY: auth-proxy-mock
auth-proxy-mock: $(APM.BIN)

# Compile a production build (excludes /dev/* and /test-fixtures/* endpoints)
$(APM.BIN): $(APM.SOURCES) Makefile
	$(GO_FLAGS) $(GO) build -o $@ cmd/auth-proxy-mock/*.go

.PHONY: apm.start-server
apm.start-server: $(APM.BIN)
	$(APM.BIN) \
		--listen-on $(APM.LISTEN_ADDR) \
		--forward-to http://$(WIP.ADDRESS)/ \
		--header "REMOTE_USER:$(APM.REMOTE_USER)" \
		--header "OIDC_CLAIM_preferred_username:$(APM.USERNAME)" \
		--header "OIDC_CLAIM_memberOf:$(APM.GROUPS)"

.PHONY: apm.start-server.usage apm.start-server.usage.current-config apm.start-server.usage.options

apm.start-server.usage: \
	apm.start-server.usage.current-config \
	apm.start-server.usage.options

apm.start-server.usage.current-config:
	@{ \
		echo ; \
		printf "current config: \n" ; \
		printf "\t make apm.start-server APM.LISTEN_ADDR=$(APM.LISTEN_ADDR) WIP.ADDRESS=$(WIP.ADDRESS) \\ \n" ; \
		printf "\t\t APM.USERNAME=$(APM.USERNAME) APM.REMOTE_USER=$(APM.REMOTE_USER) \\ \n" ; \
		printf "\t\t APM.GROUPS=$(APM.GROUPS) \n" ; \
		echo ; \
	} | $(COLORIZE)

apm.start-server.usage.options:
	@{ \
		printf ' \t MAKE OPTION \t DEFAULT VALUE \t PASSED TO \t DESCRIPTION \t\n' ; \
		printf ' \t --- \t --- \t --- \t --- \t\n' ; \
		printf ' \t `APM.LISTEN_ADDR` \t `$(APM.LISTEN_ADDR)` \t `auth-proxy-mock --listen-on=...` \t Accept requests on this address \t\n' ; \
		printf ' \t `WIP.ADDRESS` \t `$(WIP.ADDRESS)` \t `auth-proxy-mock --forward-to=...` \t Forward requests to this WiP service instance \t\n' ; \
		printf ' \t `APM.REMOTE_USER` \t `<md5sum($$USER) | as-uuid>` \t `auth-proxy-mock --header "REMOTE_USER:..."` \t Inject this sso sub in the HTTP request (header) \t\n' ; \
		printf ' \t `APM.USERNAME` \t `$$USER` \t `auth-proxy-mock --header="OIDC_CLAIM_preferred_username:..."` \t Inject this remote user in the HTTP request (header) \t\n' ; \
		printf ' \t `APM.GROUPS` \t `group1,group2,...` \t `auth-proxy-mock --header="OIDC_CLAIM_memberOf:..."` \t Inject this comma seperated group list in the HTTP request (header) \t\n' ; \
	} | column -t -s $$'\t' -o '|' | $(COLORIZE)

.PHONY: apm.docs
apm.docs: docs/auth-proxy-mock.md
docs/auth-proxy-mock.md: make.d/auth-proxy-mock.mk bin/auth-proxy-mock | /usr/bin/sed /usr/bin/sponge
	@printf "$@: inject block: <!-- make apm.start-server.usage.options -->\n"
	@{ \
		sed -n '1,/^<!-- START make apm.start-server.usage.options -->$$/ p' $@ ; \
		$(MAKE_EMBED) apm.start-server.usage.options ; \
		sed -n '/^<!-- END make apm.start-server.usage.options -->$$/,$$ p' $@ ; \
	} | sponge $@
	@
	@printf "$@: inject block: <!-- bin/auth-proxy-mock --help -->\n"
	@{ \
		sed -n '1,/^<!-- START bin\/auth-proxy-mock --help -->$$/ p' $@ ; \
		printf '```bash\n' ; \
		printf './bin/auth-proxy-mock --help\n' ; \
		./bin/auth-proxy-mock --help 2>&1 | sed 's/^/\t# /'; \
		printf '```\n' ; \
		sed -n '/^<!-- END bin\/auth-proxy-mock --help -->$$/,$$ p' $@ ; \
	} | sponge $@

	git status --short docs/auth-proxy-mock.md
	git diff docs/auth-proxy-mock.md
