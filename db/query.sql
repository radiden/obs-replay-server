-- name: UserByName :one
SELECT * FROM users WHERE name = ? LIMIT 1;

-- name: ReplaysForUser :many
SELECT * FROM replays JOIN users ON replays.owner = users.id WHERE users.name = ?;

-- name: CreateUser :one
INSERT INTO users (name) VALUES (?) RETURNING *;

-- name: AddReplay :one
INSERT INTO replays (file_path, owner) VALUES (?, ?) RETURNING *;