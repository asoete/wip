DB.SCHEMA_SOURCES := $(shell find sql/schema/ -iname '*.sql' | sort -V)

.PHONY: db db.create db.migrate db.delete

db: db.create db.migrate

db.create: | $(DB.SQLITE_FILE)

db.migrate: | $(DB.SQLITE_FILE)
	DB_FILE="$(DB.SQLITE_FILE)" libexec/db.migrate.sh

$(DB.SQLITE_FILE): sql/schema/20260825145000_create-table-migrations.sql
	sqlite3 $(DB.SQLITE_FILE) < sql/schema/20260825145000_create-table-migrations.sql

db.delete:
	@printf ">_ rm -v $(DB.SQLITE_FILE)\n"
	@read -p "are you sure (y|yes|NO): " ; \
		printf "$$REPLY" | grep -qEi '^(y|yes)$$' \
		&& { rm -v $(DB.SQLITE_FILE) || true ; } \
		|| printf ' `-> abort...\n'
