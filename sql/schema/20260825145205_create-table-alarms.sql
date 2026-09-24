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
