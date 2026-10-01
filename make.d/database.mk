DB.SCHEMA_SOURCES := $(shell find sql/schema/ -iname '*.sql' | sort -V)

.PHONY: db db.create db.migrate db.delete

db: db.create db.migrate db.sql

# ============================================================================
# FILESYSTEM
# ============================================================================

db.reset: db.delete db

DB.OPEN.MODE := -markdown
db.open: | $(DB.SQLITE_FILE)
	sqlite3 $(DB.OPEN.MODE) $(DB.SQLITE_FILE)

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
		&& { $(MAKE) db.delete.no-confirm || true ; } \
		|| printf ' `-> abort...\n'

.PHONY: db.delete.no-confirm
db.delete.no-confirm:
		rm -f $(if $(DEBUG),-v) \
			$(DB.SQLITE_FILE) \
			$(DB.SQLITE_FILE).wal \
			$(DB.SQLITE_FILE)-shm

# .PHONY: rpm.migrations
# rpm.migrations: rpm/migrate.sh

# ============================================================================
# DIST MIGTATION
# ============================================================================

rpm/migrate.sh: $(DB.SCHEMA_SOURCES) tools/bundle-migrations.sh | rpm
	tools/bundle-migrations.sh \
		$(DB.SCHEMA_SOURCES) \
		> $@

rpm:
	mkdir rpm

# ============================================================================
# SQLC CODEGEN
# ============================================================================

DB.SQLC.QUERY_SOURCES.DIR := sql/queries
DB.SQLC.QUERY_SOURCES := $(shell find $(DB.SQLC.QUERY_SOURCES.DIR) -iname '*.sql')

DB.SQLC.QUERY_TARGETS.DIR := internal/db
DB.SQLC.QUERY_TARGETS := $(addprefix $(DB.SQLC.QUERY_TARGETS.DIR)/,$(addsuffix .go,$(notdir $(DB.SQLC.QUERY_SOURCES))))

db.sql: $(DB.SQLC.QUERY_TARGETS)

$(DB.SQLC.QUERY_TARGETS.DIR)/db.go: $(DB.SCHEMA_SOURCES) $(DB.SQLC.QUERY_SOURCES) sqlc.yaml
	sqlc generate

$(DB.SQLC.QUERY_TARGETS.DIR)/models.go: $(DB.SCHEMA_SOURCES) $(DB.SQLC.QUERY_SOURCES)
	sqlc generate

$(DB.SQLC.QUERY_TARGETS.DIR)/%.go: $(DB.SQLC.QUERY_SOURCES.DIR)/%
	sqlc generate

.PHONY: db.sql.clean
db.sql.clean:
	rm -rfv $(DB.SQLC.QUERY_TARGETS)
