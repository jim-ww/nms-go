-- name: FindAllUsers :many
SELECT * FROM users;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: FindUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: InsertUser :one
 INSERT INTO users (id, username, email, password, role)
 VALUES (uuid_generate_v4(), ?, ?, ?, ?)
 RETURNING *;

-- TODO returns int64
-- name: IsUsernameTaken :one
SELECT EXISTS (
  SELECT 1 FROM users
  WHERE username = ?
);

-- name: IsEmailTaken :one
SELECT EXISTS (
  SELECT 1 FROM users
  WHERE email = ?
);



