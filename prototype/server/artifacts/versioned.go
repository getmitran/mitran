package artifacts

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const maxVersions = 50

// ArtifactStore manages versioned artifacts with slug-based identity.
type ArtifactStore struct {
	db *sql.DB
}

// Artifact represents a stored artifact with metadata.
type Artifact struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	Slug           string    `json:"slug"`
	Name           string    `json:"name"`
	Kind           string    `json:"kind"`
	Description    string    `json:"description"`
	Tags           string    `json:"tags"`
	Content        string    `json:"content"`
	CurrentVersion int       `json:"current_version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ArtifactVersion represents a single version snapshot.
type ArtifactVersion struct {
	ID         string    `json:"id"`
	ArtifactID string    `json:"artifact_id"`
	Version    int       `json:"version"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewArtifactStore creates a new store and ensures tables exist.
func NewArtifactStore(db *sql.DB) (*ArtifactStore, error) {
	store := &ArtifactStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("artifacts migration: %w", err)
	}
	return store, nil
}

func (s *ArtifactStore) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS artifacts (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			slug TEXT NOT NULL,
			name TEXT NOT NULL,
			kind TEXT NOT NULL DEFAULT 'widget',
			description TEXT DEFAULT '',
			tags TEXT DEFAULT '',
			current_version INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(project_id, slug)
		)`,
		`CREATE TABLE IF NOT EXISTS artifact_versions (
			id TEXT PRIMARY KEY,
			artifact_id TEXT NOT NULL,
			version INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (artifact_id) REFERENCES artifacts(id) ON DELETE CASCADE,
			UNIQUE(artifact_id, version)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_artifacts_project ON artifacts(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_artifact_versions_artifact ON artifact_versions(artifact_id)`,
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

// Save creates a new artifact and returns its slug.
func (s *ArtifactStore) Save(projectID, name, kind, content, tags, description string) (string, error) {
	slug := generateSlug(name)
	id := fmt.Sprintf("art_%d", time.Now().UnixNano())
	versionID := fmt.Sprintf("av_%d", time.Now().UnixNano())
	now := time.Now().UTC()

	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO artifacts (id, project_id, slug, name, kind, description, tags, current_version, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`,
		id, projectID, slug, name, kind, description, tags, now, now,
	)
	if err != nil {
		return "", fmt.Errorf("insert artifact: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO artifact_versions (id, artifact_id, version, content, created_at)
		 VALUES (?, ?, 1, ?, ?)`,
		versionID, id, content, now,
	)
	if err != nil {
		return "", fmt.Errorf("insert version: %w", err)
	}

	return slug, tx.Commit()
}

// Get retrieves an artifact. version=0 means current version.
func (s *ArtifactStore) Get(projectID, slug string, version int) (*Artifact, error) {
	var a Artifact
	err := s.db.QueryRow(
		`SELECT id, project_id, slug, name, kind, description, tags, current_version, created_at, updated_at
		 FROM artifacts WHERE project_id = ? AND slug = ?`,
		projectID, slug,
	).Scan(&a.ID, &a.ProjectID, &a.Slug, &a.Name, &a.Kind, &a.Description, &a.Tags, &a.CurrentVersion, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	targetVersion := a.CurrentVersion
	if version > 0 {
		targetVersion = version
	}

	err = s.db.QueryRow(
		`SELECT content FROM artifact_versions WHERE artifact_id = ? AND version = ?`,
		a.ID, targetVersion,
	).Scan(&a.Content)
	if err != nil {
		return nil, fmt.Errorf("version %d not found: %w", targetVersion, err)
	}

	return &a, nil
}

// Update modifies an artifact, bumps version, and prunes old versions.
func (s *ArtifactStore) Update(projectID, slug, content, name, tags, description string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id string
	var currentVersion int
	err = tx.QueryRow(
		`SELECT id, current_version FROM artifacts WHERE project_id = ? AND slug = ?`,
		projectID, slug,
	).Scan(&id, &currentVersion)
	if err != nil {
		return fmt.Errorf("artifact not found: %w", err)
	}

	newVersion := currentVersion + 1
	now := time.Now().UTC()
	versionID := fmt.Sprintf("av_%d", now.UnixNano())

	_, err = tx.Exec(
		`INSERT INTO artifact_versions (id, artifact_id, version, content, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		versionID, id, newVersion, content, now,
	)
	if err != nil {
		return fmt.Errorf("insert version: %w", err)
	}

	// Update artifact metadata
	updateQuery := `UPDATE artifacts SET current_version = ?, updated_at = ?`
	args := []interface{}{newVersion, now}

	if name != "" {
		updateQuery += `, name = ?`
		args = append(args, name)
	}
	if tags != "" {
		updateQuery += `, tags = ?`
		args = append(args, tags)
	}
	if description != "" {
		updateQuery += `, description = ?`
		args = append(args, description)
	}

	updateQuery += ` WHERE id = ?`
	args = append(args, id)

	_, err = tx.Exec(updateQuery, args...)
	if err != nil {
		return fmt.Errorf("update artifact: %w", err)
	}

	// Prune versions beyond maxVersions
	if newVersion > maxVersions {
		cutoff := newVersion - maxVersions
		_, err = tx.Exec(
			`DELETE FROM artifact_versions WHERE artifact_id = ? AND version <= ?`,
			id, cutoff,
		)
		if err != nil {
			return fmt.Errorf("prune versions: %w", err)
		}
	}

	return tx.Commit()
}

// List returns artifacts matching filters.
func (s *ArtifactStore) List(projectID, kind, tag, query string) ([]Artifact, error) {
	q := `SELECT id, project_id, slug, name, kind, description, tags, current_version, created_at, updated_at
	      FROM artifacts WHERE project_id = ?`
	args := []interface{}{projectID}

	if kind != "" {
		q += ` AND kind = ?`
		args = append(args, kind)
	}
	if tag != "" {
		q += ` AND tags LIKE ?`
		args = append(args, "%"+tag+"%")
	}
	if query != "" {
		q += ` AND name LIKE ?`
		args = append(args, "%"+query+"%")
	}

	q += ` ORDER BY updated_at DESC`

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Artifact
	for rows.Next() {
		var a Artifact
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.Slug, &a.Name, &a.Kind, &a.Description, &a.Tags, &a.CurrentVersion, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, rows.Err()
}

// Versions returns all version numbers for an artifact.
func (s *ArtifactStore) Versions(projectID, slug string) ([]int, error) {
	var id string
	err := s.db.QueryRow(
		`SELECT id FROM artifacts WHERE project_id = ? AND slug = ?`,
		projectID, slug,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(
		`SELECT version FROM artifact_versions WHERE artifact_id = ? ORDER BY version DESC`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, rows.Err()
}

// Delete removes an artifact and all its versions.
func (s *ArtifactStore) Delete(projectID, slug string) error {
	result, err := s.db.Exec(
		`DELETE FROM artifacts WHERE project_id = ? AND slug = ?`,
		projectID, slug,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("artifact %q not found in project %q", slug, projectID)
	}
	return nil
}

var slugRegex = regexp.MustCompile(`[^a-z0-9-]+`)

func generateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = slugRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	// Collapse multiple hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	if len(s) > 80 {
		s = s[:80]
		s = strings.TrimRight(s, "-")
	}
	if s == "" {
		s = fmt.Sprintf("artifact-%d", time.Now().Unix())
	}
	return s
}
