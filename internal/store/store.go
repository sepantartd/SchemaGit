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

    // Logs
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS cloud_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );
    `)

    // Projects
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS projects (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );
    `)

    return s, nil
}

//
// Cloud Logs
//

type CloudLog struct {
    ID        int64  `json:"id"`
    Type      string `json:"type"`
    Success   bool   `json:"success"`
    Error     string `json:"error"`
    Timestamp int64  `json:"timestamp"`
    ProjectID int64  `json:"project_id"`
}

func (s *Store) SaveCloudLog(l CloudLog) error {
    data, _ := json.Marshal(l)
    _, err := s.db.Exec(`INSERT INTO cloud_logs (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) LoadLogsByProject(pid int64) []CloudLog {
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC LIMIT 200`)
    if err != nil {
        return []CloudLog{}
    }

    var out []CloudLog
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        if obj.ProjectID == pid {
            out = append(out, obj)
        }
    }

    return out
}

//
// Projects
//

type Project struct {
    ID     int64  `json:"id"`
    Name   string `json:"name"`

    DBType string `json:"db_type"`

    SQLitePath string `json:"sqlite_path"`

    PGHost string `json:"pg_host"`
    PGPort int    `json:"pg_port"`
    PGUser string `json:"pg_user"`
    PGPass string `json:"pg_pass"`
    PGName string `json:"pg_name"`

    MYHost string `json:"my_host"`
    MYPort int    `json:"my_port"`
    MYUser string `json:"my_user"`
    MYPass string `json:"my_pass"`
    MYName string `json:"my_name"`
}

func (s *Store) SaveProject(p Project) error {
    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`INSERT INTO projects (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) UpdateProject(p Project) error {
    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`UPDATE projects SET data = ? WHERE id = ?`, string(data), p.ID)
    return err
}

func (s *Store) LoadProjects() []Project {
    rows, err := s.db.Query(`SELECT data FROM projects ORDER BY id ASC`)
    if err != nil {
        return []Project{}
    }

    var out []Project
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        out = append(out, obj)
    }

    return out
}
    return s, nil
}

//
// Cloud Logs
//

type CloudLog struct {
    ID        int64  `json:"id"`
    Type      string `json:"type"`
    Success   bool   `json:"success"`
    Error     string `json:"error"`
    Timestamp int64  `json:"timestamp"`
    ProjectID int64  `json:"project_id"`
}

func (s *Store) SaveCloudLog(l CloudLog) error {
    data, _ := json.Marshal(l)
    _, err := s.db.Exec(`INSERT INTO cloud_logs (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) LoadLogsByProject(pid int64) []CloudLog {
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC LIMIT 200`)
    if err != nil {
        return []CloudLog{}
    }

    var out []CloudLog
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        if obj.ProjectID == pid {
            out = append(out, obj)
        }
    }

    return out
}

//
// Projects
//

type Project struct {
    ID     int64  `json:"id"`
    Name   string `json:"name"`

    DBType string `json:"db_type"`

    SQLitePath string `json:"sqlite_path"`

    PGHost string `json:"pg_host"`
    PGPort int    `json:"pg_port"`
    PGUser string `json:"pg_user"`
    PGPass string `json:"pg_pass"`
    PGName string `json:"pg_name"`

    MYHost string `json:"my_host"`
    MYPort int    `json:"my_port"`
    MYUser string `json:"my_user"`
    MYPass string `json:"my_pass"`
    MYName string `json:"my_name"`
}

func (s *Store) SaveProject(p Project) error {
    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`INSERT INTO projects (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) UpdateProject(p Project) error {
    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`UPDATE projects SET data = ? WHERE id = ?`, string(data), p.ID)
    return err
}

func (s *Store) LoadProjects() []Project {
    rows, err := s.db.Query(`SELECT data FROM projects ORDER BY id ASC`)
    if err != nil {
        return []Project{}
    }

    var out []Project
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        out = append(out, obj)
    }

    return out
}
    return s, nil
}

//
// Cloud Logs
//

type CloudLog struct {
    ID        int64  `json:"id"`
    Type      string `json:"type"`
    Success   bool   `json:"success"`
    Error     string `json:"error"`
    Timestamp int64  `json:"timestamp"`
    ProjectID int64  `json:"project_id"`
}

func (s *Store) SaveCloudLog(l CloudLog) error {
    data, _ := json.Marshal(l)
    _, err := s.db.Exec(`INSERT INTO cloud_logs (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) LoadLogsByProject(pid int64) []CloudLog {
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC LIMIT 200`)
    if err != nil {
        return []CloudLog{}
    }

    var out []CloudLog
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        if obj.ProjectID == pid {
            out = append(out, obj)
        }
    }

    return out
}

//
// Projects
//

type Project struct {
    ID     int64  `json:"id"`
    Name   string `json:"name"`

    DBType string `json:"db_type"`

    SQLitePath string `json:"sqlite_path"`

    PGHost string `json:"pg_host"`
    PGPort int    `json:"pg_port"`
    PGUser string `json:"pg_user"`
    PGPass string `json:"pg_pass"`
    PGName string `json:"pg_name"`

    MYHost string `json:"my_host"`
    MYPort int    `json:"my_port"`
    MYUser string `json:"my_user"`
    MYPass string `json:"my_pass"`
    MYName string `json:"my_name"`
}

func (s *Store) SaveProject(p Project) error {
    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`INSERT INTO projects (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) LoadProjects() []Project {
    rows, err := s.db.Query(`SELECT data FROM projects ORDER BY id ASC`)
    if err != nil {
        return []Project{}
    }

    var out []Project
    for rows.Next() {
        var raw string
        rows.Scan(&raw)

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        out = append(out, obj)
    }

    return out
}
