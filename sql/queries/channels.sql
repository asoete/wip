-- name: InsertChannel :one
INSERT INTO ntfy_channels (name, url)
VALUES (?, ?)
RETURNING *;

-- name: SelectChannel :one
SELECT * FROM ntfy_channels WHERE name = ?;
