-- name: SelectAlarm :one
SELECT * FROM alarms WHERE alarm_id = ?;

-- name: ListUserAlarms :many
SELECT * FROM alarms WHERE user = ?;

-- name: ListActiveAlarms :many
SELECT * FROM alarms WHERE cancelled_at IS NULL;

-- name: ListActiveUserAlarms :many
SELECT * FROM alarms WHERE user = ? AND cancelled_at IS NULL ORDER BY deadline ASC;

-- name: ListCancelledUserAlarms :many
SELECT * FROM alarms WHERE user = ? AND cancelled_at IS NOT NULL ORDER BY deadline DESC;

-- name: InsertAlarm :one
INSERT INTO alarms (user, deadline, description, channel)
VALUES (?, DATETIME(?, 'utc'), ?, ?)
RETURNING *;

-- name: CancelAlarm :one
UPDATE alarms SET cancelled_at = DATETIME('now')
WHERE alarm_id = ?
RETURNING *;

-- name: UpdateAlarmDeadline :one
UPDATE alarms SET deadline = ?
WHERE alarm_id = ?
RETURNING *;
