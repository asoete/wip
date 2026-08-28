# Format the source code
.PHONY: fmt
fmt:
	@echo $(dir $(WIP_SOURCES)) | tr ' ' '\n' | sort | uniq | xargs -I {} go fmt {}

# Compile and run a DEV build
.PHONY: run
run: fmt
	WIP_CTL_TOKEN="$(CTL.TOKEN)" \
		$(GO_FLAGS) $(GO) run -tags debug cmd/wip/*.go \
			--ctl.address $(CTL.ADDRESS) \
			--web.address $(WEB.ADDRESS) \
			--db.dsn $(DB.DSN) \

# Call a Control endpoint with all necessary parameters and flags
# USAGE: make post-ctl url=/ctl/migrate/fresh
.PHONY: post-ctl
post-ctl:
	@printf ">_ POST http://$(CTL.ADDRESS)/$(url:/%=%)\n"
	@curl -i -s -X POST -H "$(CTL.TOKEN_HEADER)" \
		http://$(CTL.ADDRESS)/$(url:/%=%) \
		|& sed 's/^/|> /g'


# Kill a running server via the builtin controls
.PHONY: kill-server
kill-server:
	$(MAKE) post-ctl url=/ctl/exit/0

