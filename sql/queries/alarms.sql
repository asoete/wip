-- name: InsertAlarm :one
INSERT INTO alarms (user, deadline, description)
VALUES (?, DATETIME(?, 'utc'), ?)
RETURNING *;
