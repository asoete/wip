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
