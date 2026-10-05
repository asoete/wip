#!/usr/bin/sh

###################################################################### 
###                                                                ###
###          AUTO GENERATED DO NOT MODIFY (MANUALLY)               ###
###                                                                ###
###################################################################### 

### 
### Auto-generated via tools/bundle-migrations.sh 
### Created: ma 05 okt 2026 10:52:29 CEST
### Version: v.0.0.0-136-g5f43b91
### Commit: 5f43b9132e8cea364fef614fe5c692384d5639af
### Input-files: sql/schema/20260825145000_create-table-migrations.sql
### Input-files: sql/schema/20260825145205_create-table-alarms.sql
### Input-files: sql/schema/20260825145411_create-table-ntfy_channels.sql
### Input-files: sql/schema/20260926144043-create-table-eventlog.sql
### 

set -euo pipefail

DB_FILE="${1:-${DB_FILE:-}}"
function usage() {
	cat <<'EOUSAGE'
USAGE:
	migrate.sh <file:db.sqlite>

DESCRITPION:
	Apply all /Work in Peace/ bundled migrations to the specified SQLite
	database.
	Previously applied migrations will be skipped and only the new
	migrations will be applied (based on the `migrations` table inside the
	SQLite database).

EXAMPLE:
	$ sh dist/migrate.sh data/work-in-peace.sqlite
	sql/schema/20260825145000_create-table-migrations.sql           [  DONE  ] 2026-10-01 11:05:23
	sql/schema/20260825145205_create-table-alarms.sql               [  DONE  ] 2026-10-01 11:05:23
	sql/schema/20260825145411_create-table-ntfy_channels.sql        [  DONE  ] 2026-10-01 11:05:23
	sql/schema/20260926144043-create-table-eventlog.sql             [  DONE  ] 2026-10-01 11:05:23
EOUSAGE
}

function main() {
	if [ -z "${DB_FILE}" ] ; then
		printf " *** error: no database provided\n" 1>&2
		usage 1>&2
		exit 1
	fi

	if [ ! -f "${DB_FILE}" ] ; then
		printf " *** warning: new (uninitialized) database provided: '%s'\n" "${DB_FILE}" 1>&2
		init_new_database
	fi

	if [ ! -r "${DB_FILE}" ] ; then
		printf " *** error: invalid database provided: '%s': unreadable\n" "${DB_FILE}" 1>&2
		exit 1
	fi

	sql_schema_20260825145000_create-table-migrations_sql
	sql_schema_20260825145205_create-table-alarms_sql
	sql_schema_20260825145411_create-table-ntfy_channels_sql
	sql_schema_20260926144043-create-table-eventlog_sql
}

function sql_schema_20260825145000_create-table-migrations_sql() {

	selected_filename="$(sqlite3 "${DB_FILE}" "SELECT file FROM migrations WHERE file = 'sql/schema/20260825145000_create-table-migrations.sql';")"

	printf "%-64s" "sql/schema/20260825145000_create-table-migrations.sql"

	if [ "$selected_filename" = "sql/schema/20260825145000_create-table-migrations.sql" ] ; then
		printf "[  DONE  ] %s\n" "$(sqlite3 "${DB_FILE}" "SELECT date FROM migrations WHERE file = 'sql/schema/20260825145000_create-table-migrations.sql';")"
		return
	fi

	migrate_output=$(cat <<'EOSQL' | sqlite3 "${DB_FILE}"
CREATE TABLE IF NOT EXISTS migrations (
	migration_id INTEGER PRIMARY KEY AUTOINCREMENT,
	file TEXT NOT NULL UNIQUE,
	date TEXT NOT NULL
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145000_create-table-migrations.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
EOSQL
)

	printf "[MIGRATED] %s\n" "$migrate_output"
}

function sql_schema_20260825145205_create-table-alarms_sql() {

	selected_filename="$(sqlite3 "${DB_FILE}" "SELECT file FROM migrations WHERE file = 'sql/schema/20260825145205_create-table-alarms.sql';")"

	printf "%-64s" "sql/schema/20260825145205_create-table-alarms.sql"

	if [ "$selected_filename" = "sql/schema/20260825145205_create-table-alarms.sql" ] ; then
		printf "[  DONE  ] %s\n" "$(sqlite3 "${DB_FILE}" "SELECT date FROM migrations WHERE file = 'sql/schema/20260825145205_create-table-alarms.sql';")"
		return
	fi

	migrate_output=$(cat <<'EOSQL' | sqlite3 "${DB_FILE}"
CREATE TABLE IF NOT EXISTS alarms (
	alarm_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deadline TEXT NOT NULL,
	cancelled_at TEXT,
	channel TEXT NOT NULL,
	description TEXT
);

CREATE INDEX idx_alarms_user ON alarms(user);
CREATE INDEX idx_alarms_cancelled_at ON alarms(cancelled_at);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145205_create-table-alarms.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
EOSQL
)

	printf "[MIGRATED] %s\n" "$migrate_output"
}

function sql_schema_20260825145411_create-table-ntfy_channels_sql() {

	selected_filename="$(sqlite3 "${DB_FILE}" "SELECT file FROM migrations WHERE file = 'sql/schema/20260825145411_create-table-ntfy_channels.sql';")"

	printf "%-64s" "sql/schema/20260825145411_create-table-ntfy_channels.sql"

	if [ "$selected_filename" = "sql/schema/20260825145411_create-table-ntfy_channels.sql" ] ; then
		printf "[  DONE  ] %s\n" "$(sqlite3 "${DB_FILE}" "SELECT date FROM migrations WHERE file = 'sql/schema/20260825145411_create-table-ntfy_channels.sql';")"
		return
	fi

	migrate_output=$(cat <<'EOSQL' | sqlite3 "${DB_FILE}"
CREATE TABLE IF NOT EXISTS ntfy_channels (
	channel_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL UNIQUE
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145411_create-table-ntfy_channels.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
EOSQL
)

	printf "[MIGRATED] %s\n" "$migrate_output"
}

function sql_schema_20260926144043-create-table-eventlog_sql() {

	selected_filename="$(sqlite3 "${DB_FILE}" "SELECT file FROM migrations WHERE file = 'sql/schema/20260926144043-create-table-eventlog.sql';")"

	printf "%-64s" "sql/schema/20260926144043-create-table-eventlog.sql"

	if [ "$selected_filename" = "sql/schema/20260926144043-create-table-eventlog.sql" ] ; then
		printf "[  DONE  ] %s\n" "$(sqlite3 "${DB_FILE}" "SELECT date FROM migrations WHERE file = 'sql/schema/20260926144043-create-table-eventlog.sql';")"
		return
	fi

	migrate_output=$(cat <<'EOSQL' | sqlite3 "${DB_FILE}"
CREATE TABLE IF NOT EXISTS eventlog (
	eventlog_id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user TEXT NOT NULL,
	type TEXT NOT NULL,
	subid INTEGER,
	data TEXT
);

CREATE INDEX idx_eventlog_user ON eventlog(user);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260926144043-create-table-eventlog.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
EOSQL
)

	printf "[MIGRATED] %s\n" "$migrate_output"
}

function init_new_database() {

	printf " *** creating new database: '${DB_FILE}' and populating with 'sql/schema/20260825145000_create-table-migrations.sql'\n"
	migrate_output=$(cat <<'EOSQL' | sqlite3 "${DB_FILE}"
CREATE TABLE IF NOT EXISTS migrations (
	migration_id INTEGER PRIMARY KEY AUTOINCREMENT,
	file TEXT NOT NULL UNIQUE,
	date TEXT NOT NULL
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145000_create-table-migrations.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
EOSQL
)
	printf " *** init database '${DB_FILE}' completed\n"
	printf " *** %s\n" "$(ls -l "${DB_FILE}")" 
}
main "$@"
