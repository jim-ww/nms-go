CREATE TABLE IF NOT EXISTS notes (
      id TEXT PRIMARY KEY,
      title TEXT NOT NULL,
      content TEXT,
      user_id TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users(id)
  );

CREATE INDEX IF NOT EXISTS idx_title ON notes(title);
