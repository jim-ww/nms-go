-- name: GetUserNotesCount :one
SELECT COUNT(id) FROM notes WHERE user_id = ?;

-- name: Create :one
INSERT INTO notes(id, title, content, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ? ,?) RETURNING *;

-- 	GetByID(id int64) (*note.Note, error)
-- 	GetByTitle(title string) (*note.Note, error)
-- 	GetAll() ([]*note.Note, error)
-- 	UpdateByID(updated *note.Note) error
-- 	DeleteByID(id int64) error
-- 	DeleteByTitle(title string) error
