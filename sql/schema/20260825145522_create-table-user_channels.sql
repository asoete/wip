CREATE TABLE IF NOT EXISTS user_channels (
	uc_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	channel_id INTEGER,
	FOREIGN KEY(user_id) REFERENCES users(user_id),
	FOREIGN KEY(channel_id) REFERENCES ntfy_channels(channel_id)
);

INSERT INTO migrations (file, date)
VALUES ('sql/schema/20260825145522_create-table-user_channels.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
