// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package models

import (
	"context"
)

const addReplay = `-- name: AddReplay :one
INSERT INTO replays (file_path, owner) VALUES (?, ?) RETURNING id, file_path, owner
`

type AddReplayParams struct {
	FilePath string
	Owner    int64
}

func (q *Queries) AddReplay(ctx context.Context, arg AddReplayParams) (Replay, error) {
	row := q.db.QueryRowContext(ctx, addReplay, arg.FilePath, arg.Owner)
	var i Replay
	err := row.Scan(&i.ID, &i.FilePath, &i.Owner)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name) VALUES (?) RETURNING id, name
`

func (q *Queries) CreateUser(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, name)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const replaysForUser = `-- name: ReplaysForUser :many
SELECT replays.id, file_path, owner, users.id, name FROM replays JOIN users ON replays.owner = users.id WHERE users.name = ?
`

type ReplaysForUserRow struct {
	ID       int64
	FilePath string
	Owner    int64
	ID_2     int64
	Name     string
}

func (q *Queries) ReplaysForUser(ctx context.Context, name string) ([]ReplaysForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, replaysForUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReplaysForUserRow
	for rows.Next() {
		var i ReplaysForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.Owner,
			&i.ID_2,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const userByName = `-- name: UserByName :one
SELECT id, name FROM users WHERE name = ? LIMIT 1
`

func (q *Queries) UserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, userByName, name)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}