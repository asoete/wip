MAKE_EMBED := $(MAKE) MAKE_TERMOUT= --silent

ANSIFILTER := /usr/bin/ansifilter

COLORIZE := cat

ifneq ($(MAKE_TERMOUT),)
ifneq ($(shell which grcat 2>/dev/null), )
COLORIZE = grcat extras/grcat-colors.conf
endif
endif
