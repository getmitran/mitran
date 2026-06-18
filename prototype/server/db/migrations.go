package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	db *sql.DB
}

func Open(dir string) (*SQLiteDB, error) {
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "mitran.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	log.Printf("[db] SQLite opened at %s", path)
	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Close() error { return s.db.Close() }

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'queued',
			priority INTEGER DEFAULT 3,
			agent TEXT,
			assignee TEXT,
			labels TEXT,
			dependencies TEXT,
			result TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS tasks_fts USING fts5(title, description, result, content=tasks, content_rowid=rowid)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			languages TEXT,
			team_size INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS checkpoints (
			id TEXT PRIMARY KEY,
			task_id TEXT REFERENCES tasks(id),
			agent TEXT,
			status TEXT DEFAULT 'pending',
			description TEXT,
			files_changed TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
