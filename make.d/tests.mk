TESTS_DIR := t
TSUITE := $(TESTS_DIR)/$(suite)

bats := /usr/bin/bats
bats.flags := --pretty --timing --recursive --print-output-on-failure
make.tests.flags := --silent
ifdef DEBUG
bats.flags := --pretty --verbose-run --print-output-on-failure --show-output-of-passing-tests
make.tests.flags := 
endif

# == DATABASE ==
# ============================================================================

.PHONY: tests.newdb tests.newdb.seed

test.db.file ?= /dev/shm/wip/wip.sqlite

define tests.newdb.file.namespace-with
/dev/shm/wip/wip.$(1).sqlite
endef

tests.newdb: db.delete.no-confirm db tests.newdb.seed

tests.newdb.seed.sources := $(shell find t/db/seeds -iname '*.sql')
tests.newdb.seed:
	VERBOSITY="$(if $(DEBUG),$(DEBUG),0)" \
		DB_FILE="$(DB.SQLITE_FILE)" \
			libexec/db.migrate.sh $(tests.newdb.seed.sources)

# == TEST SUITES ==
# ============================================================================

### Note: as we use the `--` prefix to indicate "private" targets, we should
### use `make [OPTION] -- <priv-targets>` to actually execute the private
### targets...

tests.all: \
	tests.help \
	tests.config \
	tests.users \
	tests.alarms

define print-start-make-target
@printf "\n"
@printf $(if $(DEBUG),">_ ============================================================================\n","") | $(COLORIZE)
@printf ">_ *** START MAKE[%d] $@ START *** \n" "$(MAKELEVEL)" | $(COLORIZE)
@printf $(if $(DEBUG),">_ ============================================================================\n","") | $(COLORIZE)
endef

define print-end-make-target
@printf $(if $(DEBUG),">_ ============================================================================\n","") | $(COLORIZE)
@printf ">_ *** END MAKE[%d] $@ END *** \n" "$(MAKELEVEL)" | $(COLORIZE)
@printf $(if $(DEBUG),">_ ============================================================================\n","") | $(COLORIZE)
@printf "\n"
endef

tests.alarms: $(TESTS_DIR)/alarms
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags) -- \
		DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-alarms) \
		tests.env \
		tests.newdb \
		tests.seeds \
		--tests.run suite=alarms
	$(call print-end-make-target)

tests.config:
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags)  -- \
		DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-config) \
		--tests.run suite=config
	$(call print-end-make-target)

tests.env:
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags)  -- \
		DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-env) \
		--tests.run suite=env
	$(call print-end-make-target)

tests.help:
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags)  -- \
		DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-help) \
		--tests.run suite=help
	$(call print-end-make-target)

# $(intcmp $(MAKELEVEL),1, tests.newdb):
#     -> If we are summoned directly from the cli (indicated by MAKELEVEL=0),
#     run init db stuff
# See: https://www.gnu.org/software/make/manual/html_node/Conditional-Functions.html#index-intcmp
# for detail ons intcmp
tests.seeds:
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags)  -- \
		$(intcmp $(MAKELEVEL),1,DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-seeds) ) \
		$(intcmp $(MAKELEVEL),1,tests.newdb) \
		--tests.run suite=seeds
	$(call print-end-make-target)

tests.users:
	$(call print-start-make-target)
	$(MAKE) $(make.tests.flags)  -- \
		DB.SQLITE_FILE?=$(call tests.newdb.file.namespace-with,testing-users) \
		tests.newdb \
		--tests.run suite=users
	$(call print-end-make-target)

# == TEST RUNNERS ==
# ============================================================================

export DB_FILE

### Note: use the `--` prefix to indicate "private" targets
### This makes calling these targets from the cli more inconvenient
### use `make [OPTION] -- <priv-targets>` to run the private targets

.PHONY: --tests.run
--tests.run: $(WIP.BIN) | $(tests.bats)
	CTL_URL="http://$(CTL.ADDRESS)" \
	CTL_AUTH_HEADER="$(CTL.TOKEN_HEADER)" \
	PRX_URL="http://$(PRX.ADDRESS)" \
	WIP_URL="http://$(WIP.ADDRESS)" \
	DB_FILE="$(DB.SQLITE_FILE)" \
		$(bats) $(bats.flags) $(TSUITE) \
		| $(COLORIZE)

.PHONY: watch
tests.watch:
	onmod $(WIP.SOURCES) . -- 'printf "\e[H\e[22J" ; $(MAKE) tests'

.PHONY: debug
tests.debug: tests.bats.flags = --pretty --verbose-run --print-output-on-failure --show-output-of-passing-tests
tests.debug: tests.run

.PHONY: watch-debug
tests.debug.watch:
	onmod $(WIP.SOURCES) . -- 'printf "\e[H\e[22J" ; make debug '

.PHONY: %.bats
%.bats:
	bats --print-output-on-failure --timing $@

# ============================================================================
# DEV DEPENDENCIES
# ============================================================================

.PHONY: dev-dependencies
dev-dependencies: \
	/usr/bin/bats \
	/usr/bin/unbuffer\
	/usr/bin/sponge

/usr/bin/bats:
	sudo dnf install bats

/usr/bin/unbuffer:
	sudo dnf install expect

/usr/bin/sponge:
	sudo dnf install moreutils

