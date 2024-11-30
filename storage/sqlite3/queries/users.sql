-- name: FindAllUsers :many
SELECT * FROM users;

-- name: FindUserByID :one
-- SELECT * FROM users
-- WHERE id = $1;

-- name: InsertUser :one
-- INSERT INTO users (id, username, email)
-- VALUES (uuid_generate_v4(), $1, $2)
-- RETURNING *;

-- name: IsUsernameTaken :one
-- SELECT EXISTS (SELECT ONE(*) FROM users WHERE username = $1);

-- name: IsEmailTaken :one
-- SELECT EXISTS (SELECT ONE(*) FROM users WHERE email = $1);

