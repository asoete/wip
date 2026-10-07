targets := \
	/etc/systemd/system/work-in-peace.service \
	/usr/local/bin/work-in-peace \
	/usr/local/etc/work-in-peace/service.env \
	/usr/local/libexec/work-in-peace/migrate.sh \

.PHONY: list
list:
	@for file in $(targets) ; do \
		test -L "$$file" && printf "  PRESENT : " || printf "  MISSING : " ; \
		printf "$$file" ; \
		printf "\n" ; \
	done
	@printf "\n"
	@printf 'USAGE:\n'
	@printf '\tmake list (default)\n'
	@printf '\tmake install\n'

.PHONY: install
install: $(targets)

/etc/systemd/system/%: systemd/% | /etc/systemd/system/
	ln -sr $< $@

/usr/local/bin/%: bin/% | /usr/local/bin/
	ln -sr $< $@

/usr/local/etc/work-in-peace/%: etc/% | /usr/local/etc/work-in-peace/
	ln -sr $< $@

/usr/local/libexec/work-in-peace/%: libexec/% | /usr/local/libexec/work-in-peace/
	ln -sr $< $@


%/:
	mkdir -p $@
