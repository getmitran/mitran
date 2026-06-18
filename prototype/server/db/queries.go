//go:build sqlite_backend

package db

import (
	"database/sql"
	"strings"
	"time"
)

type SQLTask struct {
	ID           string
	Title        string
	Description  string
	Status       string
	Priority     int
	Agent        string
	Assignee     string
	Labels       string
	Dependencies string
	Result       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *SQLiteDB) InsertTask(t *SQLTask) error {
	_, err := s.db.Exec(`INSERT INTO tasks (id, title, description, status, priority, agent, assignee, labels, dependencies, result)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.Title, t.Description, t.Status, t.Priority, t.Agent, t.Assignee, t.Labels, t.Dependencies, t.Result)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO tasks_fts(rowid, title, description, result)
		SELECT rowid, title, description, result FROM tasks WHERE id = ?`, t.ID)
	return err
}

func (s *SQLiteDB) UpdateTask(t *SQLTask) error {
	_, err := s.db.Exec(`UPDATE tasks SET title=?, description=?, status=?, priority=?, agent=?, assignee=?, labels=?, dependencies=?, result=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		t.Title, t.Description, t.Status, t.Priority, t.Agent, t.Assignee, t.Labels, t.Dependencies, t.Result, t.ID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM tasks_fts WHERE rowid = (SELECT rowid FROM tasks WHERE id = ?)`, t.ID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO tasks_fts(rowid, title, description, result)
		SELECT rowid, title, description, result FROM tasks WHERE id = ?`, t.ID)
	return err
}

func (s *SQLiteDB) DeleteTask(id string) error {
	_, err := s.db.Exec(`DELETE FROM tasks_fts WHERE rowid = (SELECT rowid FROM tasks WHERE id = ?)`, id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func (s *SQLiteDB) GetTask(id string) (*SQLTask, error) {
	t := &SQLTask{}
	err := s.db.QueryRow(`SELECT id, title, COALESCE(description,''), status, priority, COALESCE(agent,''), COALESCE(assignee,''), COALESCE(labels,''), COALESCE(dependencies,''), COALESCE(result,''), created_at, updated_at FROM tasks WHERE id = ?`, id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Agent, &t.Assignee, &t.Labels, &t.Dependencies, &t.Result, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (s *SQLiteDB) ListTasks(status string) ([]*SQLTask, error) {
	q := `SELECT id, title, COALESCE(description,''), status, priority, COALESCE(agent,''), COALESCE(assignee,''), COALESCE(labels,''), COALESCE(dependencies,''), COALESCE(result,''), created_at, updated_at FROM tasks`
	var rows *sql.Rows
	var err error
	if status != "" {
		q += ` WHERE status = ?`
		rows, err = s.db.Query(q, status)
	} else {
		rows, err = s.db.Query(q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

func (s *SQLiteDB) SearchTasks(query string) ([]*SQLTask, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return s.ListTasks("")
	}
	rows, err := s.db.Query(`SELECT t.id, t.title, COALESCE(t.description,''), t.status, t.priority, COALESCE(t.agent,''), COALESCE(t.assignee,''), COALESCE(t.labels,''), COALESCE(t.dependencies,''), COALESCE(t.result,''), t.created_at, t.updated_at
		FROM tasks t JOIN tasks_fts f ON t.rowid = f.rowid
		WHERE tasks_fts MATCH ?
		ORDER BY rank`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*SQLTask, error) {
	var tasks []*SQLTask
	for rows.Next() {
		t := &SQLTask{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Agent, &t.Assignee, &t.Labels, &t.Dependencies, &t.Result, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
