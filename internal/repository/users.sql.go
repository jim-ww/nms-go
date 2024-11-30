// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
)

const findAllUsers = `-- name: FindAllUsers :many
SELECT id, username, email, password, role, created_at, updated_at FROM users
`

func (q *Queries) FindAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, findAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, username, email, password, role, created_at, updated_at FROM users
WHERE email = ?
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT id, username, email, password, role, created_at, updated_at FROM users
WHERE id = ?
`

func (q *Queries) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT id, username, email, password, role, created_at, updated_at FROM users
WHERE username = ?
`

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
 INSERT INTO users (id, username, email, password, role)
 VALUES (uuid_generate_v4(), ?, ?, ?, ?)
 RETURNING id, username, email, password, role, created_at, updated_at
`

type InsertUserParams struct {
	Username string
	Email    string
	Password string
	Role     role.Role
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, insertUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const isEmailTaken = `-- name: IsEmailTaken :one
SELECT EXISTS (
  SELECT 1 FROM users
  WHERE email = ?
)
`

func (q *Queries) IsEmailTaken(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, isEmailTaken, email)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const isUsernameTaken = `-- name: IsUsernameTaken :one
SELECT EXISTS (
  SELECT 1 FROM users
  WHERE username = ?
)
`

func (q *Queries) IsUsernameTaken(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, isUsernameTaken, username)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}
