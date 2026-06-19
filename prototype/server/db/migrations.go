package db

import "database/sql"

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		languages TEXT,
		team_size INTEGER,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		title TEXT,
		description TEXT,
		agent_type TEXT,
		priority INTEGER,
		status TEXT,
		depends_on TEXT,
		resources TEXT,
		created_by TEXT,
		created_at TEXT,
		started_at TEXT,
		completed_at TEXT,
		assigned_to TEXT
	);
	CREATE TABLE IF NOT EXISTS checkpoints (
		id TEXT PRIMARY KEY,
		task_id TEXT,
		agent_id TEXT,
		output TEXT,
		diffs TEXT,
		questions TEXT,
		decision_action TEXT,
		decision_feedback TEXT,
		decision_by TEXT,
		decision_at TEXT,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		name TEXT,
		type TEXT,
		status TEXT,
		callback_url TEXT,
		last_heartbeat TEXT,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS environments (
		project_id TEXT,
		name TEXT,
		auto_deploy INTEGER,
		approval_required INTEGER,
		rollback_enabled INTEGER,
		created_at TEXT,
		PRIMARY KEY (project_id, name)
	);
	CREATE TABLE IF NOT EXISTS pipelines (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		name TEXT,
		stages TEXT,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS pipeline_runs (
		id TEXT PRIMARY KEY,
		pipeline_id TEXT,
		project_id TEXT,
		branch TEXT,
		commit_hash TEXT,
		status TEXT,
		started_at TEXT,
		finished_at TEXT
	);
	CREATE TABLE IF NOT EXISTS deployments (
		id TEXT PRIMARY KEY,
		project_id TEXT,
		environment TEXT,
		version TEXT,
		commit_hash TEXT,
		status TEXT,
		triggered_at TEXT,
		finished_at TEXT
	);
	CREATE TABLE IF NOT EXISTS tickets (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		status TEXT,
		priority TEXT,
		assignee TEXT,
		labels TEXT,
		sprint_id TEXT,
		sla_due_at TEXT,
		created_by TEXT,
		created_at TEXT,
		updated_at TEXT
	);
	CREATE TABLE IF NOT EXISTS ticket_comments (
		id TEXT PRIMARY KEY,
		ticket_id TEXT,
		author TEXT,
		body TEXT,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS sprints (
		id TEXT PRIMARY KEY,
		name TEXT,
		start_date TEXT,
		end_date TEXT,
		created_at TEXT
	);
	CREATE TABLE IF NOT EXISTS wiki_pages (
		id TEXT PRIMARY KEY,
		title TEXT,
		content TEXT,
		parent_id TEXT,
		slug TEXT,
		created_by TEXT,
		tags TEXT,
		created_at TEXT,
		updated_at TEXT
	);
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT
	);
	`)
	return err
}
