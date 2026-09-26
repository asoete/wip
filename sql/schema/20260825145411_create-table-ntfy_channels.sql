CREATE TABLE IF NOT EXISTS ntfy_channels (
	channel_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL UNIQUE
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145411_create-table-ntfy_channels.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
