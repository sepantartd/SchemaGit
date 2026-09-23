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
    return initTables()
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

        CREATE TABLE IF NOT EXISTS users (
            email TEXT PRIMARY KEY,
            password TEXT
        );

        CREATE TABLE IF NOT EXISTS sessions (
            token TEXT PRIMARY KEY,
            email TEXT
        );

        CREATE TABLE IF NOT EXISTS billing (
            email TEXT PRIMARY KEY,
            plan TEXT
        );
    `)

    return err
}
