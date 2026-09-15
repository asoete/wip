-- name: SelectAlarm :one
SELECT * FROM alarms WHERE alarm_id = ?;

-- name: InsertAlarm :one
INSERT INTO alarms (user, deadline, description)
VALUES (?, DATETIME(?, 'utc'), ?)
RETURNING *;

-- name: CancelAlarm :one
UPDATE alarms SET cancelled_at = DATETIME('now')
WHERE alarm_id = ?
RETURNING *;
