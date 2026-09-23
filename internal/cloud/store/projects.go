package store

import (
    "github.com/google/uuid"
)

func CountProjects(email string) int {
    row := DB.QueryRow(`SELECT COUNT(*) FROM projects WHERE owner = ?`, email)
    var n int
    if err := row.Scan(&n); err != nil {
        return 0
    }
    return n
}

func CountMigrations(email string) int {
    row := DB.QueryRow(`SELECT COUNT(*) FROM migrations WHERE owner = ?`, email)
    var n int
    if err := row.Scan(&n); err != nil {
        return 0
    }
    return n
}

func CreateProject(email, name string) error {
    _, err := DB.Exec(`
        INSERT INTO projects (id, owner, name) VALUES (?, ?, ?)
    `, uuid.New().String(), email, name)
    return err
}

func RecordMigration(email, project string) error {
    _, err := DB.Exec(`
        INSERT INTO migrations (id, owner, project, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)
    `, uuid.New().String(), email, project)
    return err
}
