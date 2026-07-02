package memory

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// MemoryStore provides project-scoped persistent memory (facts, lessons, episodes).
type MemoryStore struct {
	db *sql.DB
}

// Fact is a key-value pair scoped to a project.
type Fact struct {
	ID        int64     `json:"id"`
	ProjectID string    `json:"project_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Lesson is a learned correction or preference.
type Lesson struct {
	ID        int64     `json:"id"`
	ProjectID string    `json:"project_id"`
	Rule      string    `json:"rule"`
	Negative  string    `json:"negative,omitempty"`
	Category  string    `json:"category"`
	Scope     string    `json:"scope"`
	CreatedAt time.Time `json:"created_at"`
}

// Episode is a timestamped memory fragment for episodic recall.
type Episode struct {
	ID            int64     `json:"id"`
	ProjectID     string    `json:"project_id"`
	Content       string    `json:"content"`
	EmbeddingText string    `json:"embedding_text,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	SessionID     string    `json:"session_id,omitempty"`
}

// NewMemoryStore creates a MemoryStore and runs migrations.
func NewMemoryStore(db *sql.DB) (*MemoryStore, error) {
	s := &MemoryStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("memory migration: %w", err)
	}
	return s, nil
}

func (s *MemoryStore) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS facts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(project_id, key)
		)`,
		`CREATE TABLE IF NOT EXISTS lessons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			rule TEXT NOT NULL,
			negative TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT 'knowledge',
			scope TEXT NOT NULL DEFAULT 'project',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS episodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			content TEXT NOT NULL,
			embedding_text TEXT NOT NULL DEFAULT '',
			timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			session_id TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX IF NOT EXISTS idx_facts_project ON facts(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_lessons_project_scope ON lessons(project_id, scope)`,
		`CREATE INDEX IF NOT EXISTS idx_episodes_project ON episodes(project_id)`,
	}
	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("exec %q: %w", m[:40], err)
		}
	}
	return nil
}

// AddFact upserts a fact (insert or update value).
func (s *MemoryStore) AddFact(projectID, key, value string) (*Fact, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO facts (project_id, key, value, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, key) DO UPDATE SET value=excluded.value, updated_at=excluded.updated_at`,
		projectID, key, value, now, now,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Fact{ID: id, ProjectID: projectID, Key: key, Value: value, CreatedAt: now, UpdatedAt: now}, nil
}

// GetFact retrieves a single fact by project and key.
func (s *MemoryStore) GetFact(projectID, key string) (*Fact, error) {
	f := &Fact{}
	err := s.db.QueryRow(
		`SELECT id, project_id, key, value, created_at, updated_at FROM facts WHERE project_id=? AND key=?`,
		projectID, key,
	).Scan(&f.ID, &f.ProjectID, &f.Key, &f.Value, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return f, err
}

// ListFacts returns all facts for a project.
func (s *MemoryStore) ListFacts(projectID string) ([]Fact, error) {
	rows, err := s.db.Query(
		`SELECT id, project_id, key, value, created_at, updated_at FROM facts WHERE project_id=? ORDER BY key`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var facts []Fact
	for rows.Next() {
		var f Fact
		if err := rows.Scan(&f.ID, &f.ProjectID, &f.Key, &f.Value, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		facts = append(facts, f)
	}
	return facts, rows.Err()
}

// DeleteFact removes a fact by project and key.
func (s *MemoryStore) DeleteFact(projectID, key string) error {
	_, err := s.db.Exec(`DELETE FROM facts WHERE project_id=? AND key=?`, projectID, key)
	return err
}

// AddLesson stores a new lesson.
func (s *MemoryStore) AddLesson(projectID, rule, negative, category, scope string) (*Lesson, error) {
	if category == "" {
		category = "knowledge"
	}
	if scope == "" {
		scope = "project"
	}
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO lessons (project_id, rule, negative, category, scope, created_at) VALUES (?,?,?,?,?,?)`,
		projectID, rule, negative, category, scope, now,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Lesson{ID: id, ProjectID: projectID, Rule: rule, Negative: negative, Category: category, Scope: scope, CreatedAt: now}, nil
}

// ListLessons returns lessons for a project, optionally filtered by scope.
func (s *MemoryStore) ListLessons(projectID, scope string) ([]Lesson, error) {
	var rows *sql.Rows
	var err error
	if scope != "" {
		rows, err = s.db.Query(
			`SELECT id, project_id, rule, negative, category, scope, created_at FROM lessons WHERE project_id=? AND scope=? ORDER BY created_at DESC`,
			projectID, scope,
		)
	} else {
		rows, err = s.db.Query(
			`SELECT id, project_id, rule, negative, category, scope, created_at FROM lessons WHERE project_id=? ORDER BY created_at DESC`,
			projectID,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.Rule, &l.Negative, &l.Category, &l.Scope, &l.CreatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, rows.Err()
}

// RemoveLesson deletes lessons whose rule contains the query substring.
func (s *MemoryStore) RemoveLesson(projectID, query string) (int64, error) {
	res, err := s.db.Exec(
		`DELETE FROM lessons WHERE project_id=? AND rule LIKE ?`,
		projectID, "%"+query+"%",
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// AddEpisode stores a new episodic memory fragment.
func (s *MemoryStore) AddEpisode(projectID, content, sessionID string) (*Episode, error) {
	now := time.Now().UTC()
	// embedding_text is content lowercased for LIKE search; vector search in v0.2.0
	embeddingText := strings.ToLower(content)
	res, err := s.db.Exec(
		`INSERT INTO episodes (project_id, content, embedding_text, timestamp, session_id) VALUES (?,?,?,?,?)`,
		projectID, content, embeddingText, now, sessionID,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Episode{ID: id, ProjectID: projectID, Content: content, EmbeddingText: embeddingText, Timestamp: now, SessionID: sessionID}, nil
}

// SearchEpisodes performs LIKE-based search on episodes. Vector search deferred to v0.2.0.
func (s *MemoryStore) SearchEpisodes(projectID, query string, limit int) ([]Episode, error) {
	if limit <= 0 {
		limit = 10
	}
	q := "%" + strings.ToLower(query) + "%"
	rows, err := s.db.Query(
		`SELECT id, project_id, content, embedding_text, timestamp, session_id
		 FROM episodes WHERE project_id=? AND embedding_text LIKE ?
		 ORDER BY timestamp DESC LIMIT ?`,
		projectID, q, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var episodes []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Content, &e.EmbeddingText, &e.Timestamp, &e.SessionID); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}
