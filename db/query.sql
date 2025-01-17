-- name: UserByName :one
SELECT * FROM users WHERE name = ? LIMIT 1;

-- name: ReplaysForUser :many
SELECT * FROM replays JOIN users ON replays.owner = users.id WHERE users.name = ? ORDER BY replays.id DESC;

-- name: CreateUser :one
INSERT INTO users (name) VALUES (?) RETURNING *;

-- name: AddReplay :one
INSERT INTO replays (file_path, creation_time, owner) VALUES (?, ?, ?) RETURNING *;