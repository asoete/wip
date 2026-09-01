TESTS_DIR := t
TESTS := $(TESTS_DIR)/$(T)

.PHONY: tests
tests:
	bats --timing --recursive --print-output-on-failure $(TESTS)

.PHONY: watch
tests.watch:
	onmod $(WIP.SOURCES) . -- 'printf "\e[H\e[22J" ; $(MAKE) tests'

.PHONY: debug
tests.debug:
	$(MAKE) -C $(WIP.SOURCES)
	bats --verbose-run --print-output-on-failure --show-output-of-passing-tests $(TESTS)

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

