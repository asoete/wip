CREATE TABLE IF NOT EXISTS alarms (
	alarm_id INTEGER PRIMARY KEY AUTOINCREMENT,
	owner_id INTEGER,
	created_at TEXT NOT NULL DEFAULT 'now',
	deadline TEXT NOT NULL,
	cancelled_at TEXT,
	description TEXT,
	FOREIGN KEY(owner_id) REFERENCES users(user_id)
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145205_create-table-alarms.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
