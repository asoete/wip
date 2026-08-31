CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	sso_sub TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145030_create-table-users.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
