CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT(30) NOT NULL UNIQUE,
			email TEXT(255) NOT NULL UNIQUE,
			password TEXT(255) NOT NULL,
			role TEXT(10) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
	);

CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON users(username);

CREATE INDEX IF NOT EXISTS idx_email ON users(email);