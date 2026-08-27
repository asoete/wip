DB.SCHEMA_SOURCES := $(shell find sql/schema/ -iname '*.sql' | sort -V)

.PHONY: db db.create db.migrate db.delete

db: db.create db.migrate

db.create: | $(DB.SQLITE_FILE)

db.migrate: | $(DB.SQLITE_FILE)
	libexec/db.migrate.sh

$(DB.SQLITE_FILE): sql/schema/20260825145000_create-table-migrations.sql
	sqlite3 $(DB.SQLITE_FILE) < sql/schema/20260825145000_create-table-migrations.sql

db.delete:
	rm -v $(DB.SQLITE_FILE) || true
