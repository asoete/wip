TESTS_DIR := t
TESTS := $(TESTS_DIR)/$(T)

tests.bats := /usr/bin/bats
tests.bats.flags := --pretty --timing --recursive --print-output-on-failure

.PHONY: tests
tests: $(WIP.BIN) | $(tests.bats)
	$(tests.bats) $(tests.bats.flags) $(TESTS) | $(COLORIZE)

.PHONY: watch
tests.watch:
	onmod $(WIP.SOURCES) . -- 'printf "\e[H\e[22J" ; $(MAKE) tests'

.PHONY: debug
tests.debug: tests.bats.flags = --pretty --verbose-run --print-output-on-failure --show-output-of-passing-tests
tests.debug: tests

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

