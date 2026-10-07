.PHONY: dist
dist: \
	dist/libexec/migrate.sh \
	dist/systemd/work-in-peace.service \
	dist/etc/service.env \
	dist/bin/work-in-peace

dist/bin/work-in-peace: bin/wip | dist/bin/
	cp -av $< $@

dist/libexec/migrate.sh: $(DB.SCHEMA_SOURCES) tools/bundle-migrations.sh | dist/libexec/
	tools/bundle-migrations.sh \
		$(DB.SCHEMA_SOURCES) \
		> $@
	chmod a+x $@

dist/etc/service.env: dist/bin/work-in-peace make.d/dist.mk | dist/etc/
	@printf '# The SQLite database file\n' > $@
	@printf '# This is the only required config and is passed to:\n' >> $@
	@printf '#  * work-in-peace.service:ExecPre=.../migrate.sh $${DB_FILE}\n' >> $@
	@printf '#  * work-in-peace.service:ExecStart=work-in-peace --db.dsn=file:$${DB_FILE}\n' >> $@
	@printf '# work-in-peace.service will FAIL if DB_FILE is empty / unset\n' >> $@
	@printf 'DB_FILE=""\n' >> $@
	@printf '\n' >> $@
	@printf '# OPTIONS= \\\n' >> $@
	@printf '#\t--web.address=127.0.0.1:8080 \\\n' >> $@
	@printf '#\t [...more options...] \\\n' >> $@
	@printf '\n' >> $@
	@printf '# Available Options:\n' >> $@
	@$< --help 2>&1 | tail -n +2 | sed 's/^/# /'>> $@

dist/bin/ dist/libexec/ dist/etc/ dist/systemd/ :
	mkdir -p $@

.PHONY: dist.clean
dist.clean:
	rm -rfv dist/{libexec,etc,bin}/*
