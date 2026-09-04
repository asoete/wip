INSERT INTO ntfy_channels (url) VALUES
('https://ntfy.sh/WiP-testing-channel_OfTD36Z4uARQeJsQB76pNQMoloPv');

INSERT INTO migrations (file, date)
VALUES ('t/db/seeds/insert-ntfy_channels.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
