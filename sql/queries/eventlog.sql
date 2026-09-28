-- name: InsertEvent :exec
INSERT INTO eventlog
	(timestamp, user, type, subid, data)
VALUES
	(DATETIME('NOW'), ?, ?, ?, ?);

