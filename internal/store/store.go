package store

import (
    "database/sql"
    "encoding/json"
    "os"

    "github.com/sepanta/schemagit/internal/schema"
    _ "github.com/mattn/go-sqlite3"
)

type Store struct {
    db *sql.DB
}

func Open(path string) (*Store, error) {
    _, err := os.Stat(path)
    first := os.IsNotExist(err)

    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    s := &Store{db: db}

    if first {
        if err := s.init(); err != nil {
            return nil, err
        }
    }

    return s, nil
}

func (s *Store) init() error {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_history (
            commit_hash TEXT PRIMARY KEY,
            schema_json TEXT NOT NULL
        );
    `)
    return err
}

func (s *Store) SaveSchema(commit string, sc *schema.DatabaseSchema) error {
    data, err := json.Marshal(sc)
    if err != nil {
        return err
    }

    _, err = s.db.Exec(
        "INSERT OR REPLACE INTO schema_history (commit_hash, schema_json) VALUES (?, ?)",
        commit, string(data),
    )
    return err
}

func (s *Store) LoadSchema(commit string) (*schema.DatabaseSchema, error) {
    row := s.db.QueryRow(
        "SELECT schema_json FROM schema_history WHERE commit_hash = ?",
        commit,
    )

    var jsonStr string
    if err := row.Scan(&jsonStr); err != nil {
        return nil, err
    }

    var sc schema.DatabaseSchema
    if err := json.Unmarshal([]byte(jsonStr), &sc); err != nil {
        return nil, err
    }

    return &sc, nil
}
