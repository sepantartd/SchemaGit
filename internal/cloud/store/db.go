package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	DB = db
	if err := initTables(); err != nil {
		db.Close()
		DB = nil
		return err
	}
	return nil
}

func initTables() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS plans (
			id TEXT PRIMARY KEY,
			repo TEXT,
			pr_number INTEGER,
			plan_json TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			status TEXT,
			result TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS logs (
			id TEXT PRIMARY KEY,
			job_id TEXT,
			log TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS api_tokens (token TEXT PRIMARY KEY, email TEXT);
		CREATE TABLE IF NOT EXISTS project_schema (owner TEXT, project TEXT, schema TEXT);
		CREATE TABLE IF NOT EXISTS users (email TEXT PRIMARY KEY, password TEXT);
		CREATE TABLE IF NOT EXISTS sessions (token TEXT PRIMARY KEY, email TEXT);
		CREATE TABLE IF NOT EXISTS orgs (id TEXT PRIMARY KEY, name TEXT, owner TEXT);
		CREATE TABLE IF NOT EXISTS org_members (org_id TEXT, email TEXT, role TEXT,
			UNIQUE(org_id, email));
		CREATE TABLE IF NOT EXISTS webhooks (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			url TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_webhooks_org_id ON webhooks(org_id);
		CREATE TABLE IF NOT EXISTS notifications (
			org_id TEXT PRIMARY KEY,
			email BOOLEAN,
			slack BOOLEAN,
			slack_webhook TEXT
		);
		CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			org_id TEXT,
			event TEXT,
			detail TEXT,
			created_at TEXT
		);
	`)
	return err
}
