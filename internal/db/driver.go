package db

import "github.com/sepanta/schemagit/internal/schema"

type DBType string

const (
    SQLite   DBType = "sqlite"
    Postgres DBType = "postgres"
    MySQL    DBType = "mysql"
)

type Driver interface {
    LoadSchema() (*schema.DatabaseSchema, error)
    ApplySQL([]string) error
}
