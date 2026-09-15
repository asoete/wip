-- name: SelectAlarm :one
SELECT * FROM alarms WHERE alarm_id = ?;

-- name: ListUserAlarms :many
SELECT * FROM alarms WHERE user = ?;

-- name: ListActiveUserAlarms :many
SELECT * FROM alarms WHERE user = ? AND cancelled_at IS NULL ORDER BY deadline ASC;

-- name: ListCancelledUserAlarms :many
SELECT * FROM alarms WHERE user = ? AND cancelled_at IS NOT NULL;

-- name: InsertAlarm :one
INSERT INTO alarms (user, deadline, description)
VALUES (?, DATETIME(?, 'utc'), ?)
RETURNING *;

-- name: CancelAlarm :one
UPDATE alarms SET cancelled_at = DATETIME('now')
WHERE alarm_id = ?
RETURNING *;
