DB.SCHEMA_SOURCES := $(shell find sql/schema/ -iname '*.sql' | sort -V)

.PHONY: db db.create db.migrate db.delete

db: db.create db.migrate

db.create: | $(DB.SQLITE_FILE)

db.migrate: | $(DB.SQLITE_FILE)
	VERBOSITY="$(if $(DEBUG),$(DEBUG),0)" \
		DB_FILE="$(DB.SQLITE_FILE)" libexec/db.migrate.sh

$(DB.SQLITE_FILE): sql/schema/20260825145000_create-table-migrations.sql | $(dir $(DB.SQLITE_FILE))
	sqlite3 $(DB.SQLITE_FILE) < sql/schema/20260825145000_create-table-migrations.sql >/dev/null

$(dir $(DB.SQLITE_FILE)):
	mkdir -p $@

db.delete:
	@printf ">_ rm -v $(DB.SQLITE_FILE)\n"
	@read -p "are you sure (y|yes|NO): " ; \
		printf "$$REPLY" | grep -qEi '^(y|yes)$$' \
		&& { $(MAKE db.delete.no-confirm) || true ; } \
		|| printf ' `-> abort...\n'

.PHONY: db.delete.no-confirm
db.delete.no-confirm:
		rm -f $(if $(DEBUG),-v) $(DB.SQLITE_FILE)

# .PHONY: rpm.migrations
# rpm.migrations: rpm/migrate.sh

rpm/migrate.sh: $(DB.SCHEMA_SOURCES) tools/bundle-migrations.sh | rpm
	tools/bundle-migrations.sh \
		$(DB.SCHEMA_SOURCES) \
		> $@

rpm:
	mkdir rpm
