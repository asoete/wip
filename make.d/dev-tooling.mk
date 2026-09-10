# Format the source code
.PHONY: fmt
fmt:
	@echo $(dir $(WIP.SOURCES)) | tr ' ' '\n' | sort | uniq | xargs -I {} go fmt {}

# Compile and run a DEV build
.PHONY: run
run: fmt
	WIP_CTL_TOKEN="$(CTL.TOKEN)" \
		$(GO_FLAGS) $(GO) run -tags debug cmd/wip/*.go \
			--ctl.address $(CTL.ADDRESS) \
			--web.address $(WIP.ADDRESS) \
			--pidfile $(WIP.PIDFILE) \
			--db.dsn $(DB.DSN) \

# Call a Control endpoint with all necessary parameters and flags
# USAGE: make post-ctl url=/ctl/migrate/fresh
.PHONY: post-ctl
post-ctl:
	@printf ">_ POST http://$(CTL.ADDRESS)/$(url:/%=%)\n"
	@curl -i -s -X POST -H "$(CTL.TOKEN_HEADER)" \
		http://$(CTL.ADDRESS)/$(url:/%=%) \
		|& sed 's/^/|> /g'

# USAGE: make get-ctl url=/ctl/migrate/fresh
.PHONY: get-ctl
get-ctl:
	@printf ">_ GET http://$(CTL.ADDRESS)/$(url:/%=%)\n"
	@curl -i -s -X GET -H "$(CTL.TOKEN_HEADER)" \
		http://$(CTL.ADDRESS)/$(url:/%=%) \
		|& sed 's/^/|> /g'

# Kill a running server via the builtin controls
.PHONY: web.server.kill
wip.server.kill:
	$(MAKE) post-ctl url=/ctl/exit/0

# Start a wip.server loop: auto restart a wip.server when it exists
.PHONY: wip.server-loop.start
wip.server-loop.start: /usr/bin/inotifywait
	while true ; do \
		printf "\e[H\e[22J" ; \
		$(MAKE) --no-print-directory run || ONESHOT=1 ./tools/notifywait.sh ': done waiting for bugfix' ; \
		sleep 0.1 ; \
	done

# Start a watcher which will kill a running wip.server on source file modifications
.PHONY: wip.server-killer.start
wip.server-killer.start: | tools/notifywait.sh
	./tools/notifywait.sh \
		'$(MAKE) --no-print-directory wip.server.kill' \
		$(shell find * -maxdepth 0 -not -iname 'bin')

# Start all required services for "a good development experience" (TM)
.PHONY: serve
serve:
	$(MAKE) --no-print-directory -j \
		wip.server-loop.start \
		wip.server-killer.start \
		apm.start-server \
		|& $(COLORIZE)
