package store

import (
    "database/sql"
    "encoding/json"

    _ "github.com/mattn/go-sqlite3"
)

type Store struct {
    db *sql.DB
}

func Open(path string) (*Store, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    s := &Store{db: db}

    // init tables
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS cloud_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        )
    `)

    return s, nil
}

func (s *Store) SaveCloudLog(l interface{}) error {
    data, _ := json.Marshal(l)
    _, err := s.db.Exec(`INSERT INTO cloud_logs (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) LoadCloudLogs() []CloudLog {
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC LIMIT 100`)
    if err != nil {
        return []CloudLog{}
    }

    var out []CloudLog

    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        out = append(out, obj)
    }

    return out
}
