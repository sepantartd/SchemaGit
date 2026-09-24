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

func UserProjects(email string) []string {
    rows, _ := DB.Query(`
        SELECT name FROM projects WHERE owner = ?
    `, email)

    var out []string
    for rows.Next() {
        var name string
        rows.Scan(&name)
        out = append(out, name)
    }
    return out
}

func LoadProjectSchema(email, project string) (string, error) {
    row := DB.QueryRow(`
        SELECT schema FROM project_schema
        WHERE owner = ? AND project = ?
    `, email, project)

    var schema string
    err := row.Scan(&schema)
    return schema, err
}

func LoadProjectMigrations(email, project string) ([]string, error) {
    rows, err := DB.Query(`
        SELECT id FROM migrations
        WHERE owner = ? AND project = ?
        ORDER BY created_at ASC
    `, email, project)
    if err != nil {
        return nil, err
    }

    var out []string
    for rows.Next() {
        var id string
        rows.Scan(&id)
        out = append(out, id)
    }
    return out, nil
}
