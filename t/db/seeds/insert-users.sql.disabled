INSERT INTO users (username, sso_sub) VALUES
('johnd', '0db52b4e-a61f-fc3f-58b4-d21c237151a1'),
('janed', 'd2f9f861-496c-4639-6a1c-666a5a781406');

INSERT INTO migrations (file, date)
VALUES ('t/db/seeds/insert-users.sql', datetime('now'))
RETURNING datetime(date, 'localtime');
