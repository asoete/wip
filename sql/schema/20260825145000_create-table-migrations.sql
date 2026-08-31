CREATE TABLE IF NOT EXISTS migrations (
	migration_id INTEGER PRIMARY KEY AUTOINCREMENT,
	file TEXT NOT NULL UNIQUE,
	date TEXT NOT NULL
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145000_create-table-migrations.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
