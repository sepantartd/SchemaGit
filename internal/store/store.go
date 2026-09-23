package store

import (
    "database/sql"
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "time"

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

    // Cloud Logs
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS cloud_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );
    `)

    // Cloud Projects
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS projects (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );
    `)

    // Project Members
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS project_members (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            project_id INTEGER NOT NULL,
            email TEXT NOT NULL,
            role TEXT NOT NULL,
            invited_at INTEGER NOT NULL
        );
    `)

    // Project Webhooks
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS project_webhooks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            project_id INTEGER NOT NULL,
            url TEXT NOT NULL,
            secret TEXT NOT NULL,
            type TEXT NOT NULL,
            created_at INTEGER NOT NULL
        );
    `)

    // Jobs Queue
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS jobs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT NOT NULL,
            payload TEXT NOT NULL,
            status TEXT NOT NULL,
            created_at INTEGER NOT NULL,
            updated_at INTEGER NOT NULL,
            attempts INTEGER NOT NULL,
            last_error TEXT
        );
    `)

    // Agent Logs
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS agent_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );
    `)

    // Migrations
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            applied_at INTEGER NOT NULL,
            checksum TEXT NOT NULL,
            direction TEXT NOT NULL
        );
    `)

    // Schema Snapshots
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_snapshots (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            created_at INTEGER NOT NULL,
            data TEXT NOT NULL
        );
    `)

    // Local Cache
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS local_cache (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL,
            updated_at INTEGER NOT NULL
        );
    `)

    // Settings
    _, _ = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS settings (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL
        );
    `)

    return s, nil
}

func (s *Store) Close() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

//
// -----------------------------------------------------------------------------
// Cloud Logs
// -----------------------------------------------------------------------------
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
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC LIMIT 500`)
    if err != nil {
        return []CloudLog{}
    }
    defer rows.Close()

    var out []CloudLog
    for rows.Next() {
        var raw string
        if err := rows.Scan(&raw); err != nil {
            continue
        }

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        if obj.ProjectID == pid {
            out = append(out, obj)
        }
    }

    return out
}

func (s *Store) DeleteLogsByProject(pid int64) error {
    rows, err := s.db.Query(`SELECT rowid, data FROM cloud_logs`)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var rowid int64
        var raw string
        if err := rows.Scan(&rowid, &raw); err != nil {
            return err
        }

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)

        if obj.ProjectID == pid {
            _, _ = s.db.Exec(`DELETE FROM cloud_logs WHERE rowid = ?`, rowid)
        }
    }

    return nil
}

func (s *Store) LoadAllLogs() []CloudLog {
    rows, err := s.db.Query(`SELECT data FROM cloud_logs ORDER BY id DESC`)
    if err != nil {
        return []CloudLog{}
    }
    defer rows.Close()

    var out []CloudLog
    for rows.Next() {
        var raw string
        if err := rows.Scan(&raw); err != nil {
            continue
        }

        var obj CloudLog
        json.Unmarshal([]byte(raw), &obj)
        out = append(out, obj)
    }

    return out
}

//
// -----------------------------------------------------------------------------
// Projects
// -----------------------------------------------------------------------------
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

    GitHubWebhookURL string `json:"github_webhook_url"`
    GitHubSecret     string `json:"github_secret"`

    CreatedAt int64 `json:"created_at"`
    UpdatedAt int64 `json:"updated_at"`
}

func (s *Store) SaveProject(p Project) error {
    if p.CreatedAt == 0 {
        p.CreatedAt = time.Now().Unix()
    }
    p.UpdatedAt = time.Now().Unix()

    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`INSERT INTO projects (data) VALUES (?)`, string(data))
    return err
}

func (s *Store) UpdateProject(p Project) error {
    p.UpdatedAt = time.Now().Unix()

    data, _ := json.Marshal(p)
    _, err := s.db.Exec(`UPDATE projects SET data = ? WHERE id = ?`, string(data), p.ID)
    return err
}

func (s *Store) LoadProjects() []Project {
    rows, err := s.db.Query(`SELECT rowid, data FROM projects ORDER BY rowid ASC`)
    if err != nil {
        return []Project{}
    }
    defer rows.Close()

    var out []Project
    idx := int64(1)

    for rows.Next() {
        var rowid int64
        var raw string
        if err := rows.Scan(&rowid, &raw); err != nil {
            continue
        }

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        if obj.ID == 0 {
            obj.ID = idx
        }

        out = append(out, obj)
        idx++
    }

    return out
}

func (s *Store) DeleteProject(id int64) error {
    rows, err := s.db.Query(`SELECT rowid, data FROM projects`)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var rowid int64
        var raw string
        if err := rows.Scan(&rowid, &raw); err != nil {
            return err
        }

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        if obj.ID == id {
            _, _ = s.db.Exec(`DELETE FROM projects WHERE rowid = ?`, rowid)
        }
    }

    return nil
}

func (s *Store) GetProjectByName(name string) (*Project, error) {
    rows, err := s.db.Query(`SELECT data FROM projects`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var raw string
        if err := rows.Scan(&raw); err != nil {
            continue
        }

        var obj Project
        json.Unmarshal([]byte(raw), &obj)

        if obj.Name == name {
            return &obj, nil
        }
    }

    return nil, errors.New("project not found")
}

//
// -----------------------------------------------------------------------------
// Project Members
// -----------------------------------------------------------------------------
//

type ProjectMember struct {
    ID        int64  `json:"id"`
    ProjectID int64  `json:"project_id"`
    Email     string `json:"email"`
    Role      string `json:"role"` // owner, admin, write, read
    InvitedAt int64  `json:"invited_at"`
}

func (s *Store) AddProjectMember(m ProjectMember) error {
    if m.InvitedAt == 0 {
        m.InvitedAt = time.Now().Unix()
    }
    _, err := s.db.Exec(`
        INSERT INTO project_members (project_id, email, role, invited_at)
        VALUES (?, ?, ?, ?)
    `, m.ProjectID, m.Email, m.Role, m.InvitedAt)
    return err
}

func (s *Store) LoadMembersByProject(pid int64) []ProjectMember {
    rows, err := s.db.Query(`
        SELECT id, project_id, email, role, invited_at
        FROM project_members
        WHERE project_id = ?
        ORDER BY invited_at ASC
    `, pid)
    if err != nil {
        return []ProjectMember{}
    }
    defer rows.Close()

    var out []ProjectMember
    for rows.Next() {
        var m ProjectMember
        if err := rows.Scan(&m.ID, &m.ProjectID, &m.Email, &m.Role, &m.InvitedAt); err != nil {
            continue
        }
        out = append(out, m)
    }

    return out
}

func (s *Store) DeleteMember(id int64) error {
    _, err := s.db.Exec(`DELETE FROM project_members WHERE id = ?`, id)
    return err
}

func (s *Store) DeleteMembersByProject(pid int64) error {
    _, err := s.db.Exec(`DELETE FROM project_members WHERE project_id = ?`, pid)
    return err
}

//
// -----------------------------------------------------------------------------
// Project Webhooks
// -----------------------------------------------------------------------------
//

type ProjectWebhook struct {
    ID        int64  `json:"id"`
    ProjectID int64  `json:"project_id"`
    URL       string `json:"url"`
    Secret    string `json:"secret"`
    Type      string `json:"type"` // github, gitlab, bitbucket, custom
    CreatedAt int64  `json:"created_at"`
}

func (s *Store) AddWebhook(w ProjectWebhook) error {
    if w.CreatedAt == 0 {
        w.CreatedAt = time.Now().Unix()
    }
    _, err := s.db.Exec(`
        INSERT INTO project_webhooks (project_id, url, secret, type, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, w.ProjectID, w.URL, w.Secret, w.Type, w.CreatedAt)
    return err
}

func (s *Store) LoadWebhooksByProject(pid int64) []ProjectWebhook {
    rows, err := s.db.Query(`
        SELECT id, project_id, url, secret, type, created_at
        FROM project_webhooks
        WHERE project_id = ?
        ORDER BY created_at ASC
    `, pid)
    if err != nil {
        return []ProjectWebhook{}
    }
    defer rows.Close()

    var out []ProjectWebhook
    for rows.Next() {
        var w ProjectWebhook
        if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &w.Secret, &w.Type, &w.CreatedAt); err != nil {
            continue
        }
        out = append(out, w)
    }

    return out
}

func (s *Store) DeleteWebhook(id int64) error {
    _, err := s.db.Exec(`DELETE FROM project_webhooks WHERE id = ?`, id)
    return err
}

func (s *Store) DeleteWebhooksByProject(pid int64) error {
    _, err := s.db.Exec(`DELETE FROM project_webhooks WHERE project_id = ?`, pid)
    return err
}

//
// -----------------------------------------------------------------------------
// Jobs Queue
// -----------------------------------------------------------------------------
//

type Job struct {
    ID        int64             `json:"id"`
    Type      string            `json:"type"`
    Payload   map[string]string `json:"payload"`
    Status    string            `json:"status"` // pending, running, done, failed
    CreatedAt int64             `json:"created_at"`
    UpdatedAt int64             `json:"updated_at"`
    Attempts  int               `json:"attempts"`
    LastError string            `json:"last_error"`
}

func (s *Store) CreateJob(j Job) error {
    raw, _ := json.Marshal(j.Payload)
    now := time.Now().Unix()

    _, err := s.db.Exec(`
        INSERT INTO jobs (type, payload, status, created_at, updated_at, attempts)
        VALUES (?, ?, 'pending', ?, ?, 0)
    `, j.Type, string(raw), now, now)

    return err
}

func (s *Store) FetchNextPendingJob() (*Job, error) {
    row := s.db.QueryRow(`
        SELECT id, type, payload, status, created_at, updated_at, attempts, last_error
        FROM jobs
        WHERE status = 'pending'
        ORDER BY id ASC
        LIMIT 1
    `)

    var (
        id, created, updated int64
        attempts             int
        typ, payloadRaw, status, lastErr string
    )

    err := row.Scan(&id, &typ, &payloadRaw, &status, &created, &updated, &attempts, &lastErr)
    if err != nil {
        return nil, err
    }

    var payload map[string]string
    json.Unmarshal([]byte(payloadRaw), &payload)

    return &Job{
        ID:        id,
        Type:      typ,
        Payload:   payload,
        Status:    status,
        CreatedAt: created,
        UpdatedAt: updated,
        Attempts:  attempts,
        LastError: lastErr,
    }, nil
}

func (s *Store) UpdateJobStatus(id int64, status string, errMsg string) {
    now := time.Now().Unix()
    _, _ = s.db.Exec(`
        UPDATE jobs
        SET status = ?, updated_at = ?, last_error = ?
        WHERE id = ?
    `, status, now, errMsg, id)
}

func (s *Store) IncrementJobAttempts(id int64) {
    _, _ = s.db.Exec(`
        UPDATE jobs
        SET attempts = attempts + 1
        WHERE id = ?
    `, id)
}

func (s *Store) LoadJobHistory() []Job {
    rows, err := s.db.Query(`
        SELECT id, type, payload, status, created_at, updated_at, attempts, last_error
        FROM jobs
        ORDER BY id DESC
        LIMIT 200
    `)
    if err != nil {
        return []Job{}
    }
    defer rows.Close()

    var out []Job
    for rows.Next() {
        var (
            id, created, updated int64
            attempts             int
            typ, payloadRaw, status, lastErr string
        )

        if err := rows.Scan(&id, &typ, &payloadRaw, &status, &created, &updated, &attempts, &lastErr); err != nil {
            continue
        }

        var payload map[string]string
        json.Unmarshal([]byte(payloadRaw), &payload)

        out = append(out, Job{
            ID:        id,
            Type:      typ,
            Payload:   payload,
            Status:    status,
            CreatedAt: created,
            UpdatedAt: updated,
            Attempts:  attempts,
            LastError: lastErr,
        })
    }

    return out
}

//
// -----------------------------------------------------------------------------
// Agent Logs
// -----------------------------------------------------------------------------
//

type AgentLog struct {
    ID        int64  `json:"id"`
    JobID     int64  `json:"job_id"`
    ProjectID int64  `json:"project_id"`
    Type      string `json:"type"`     // push, diff, apply, error, retry
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}

func (s *Store) SaveAgentLog(l AgentLog) error {
    raw, _ := json.Marshal(l)
    _, err := s.db.Exec(`
        INSERT INTO agent_logs (data)
        VALUES (?)
    `, string(raw))
    return err
}

func (s *Store) LoadAgentLogsByProject(pid int64) []AgentLog {
    rows, err := s.db.Query(`
        SELECT data FROM agent_logs ORDER BY id DESC LIMIT 500
    `)
    if err != nil {
        return []AgentLog{}
    }
    defer rows.Close()

    var out []AgentLog
    for rows.Next() {
        var raw string
        if err := rows.Scan(&raw); err != nil {
            continue
        }

        var obj AgentLog
        if json.Unmarshal([]byte(raw), &obj) == nil && obj.ProjectID == pid {
            out = append(out, obj)
        }
    }

    return out
}

//
// -----------------------------------------------------------------------------
// Migrations Store
// -----------------------------------------------------------------------------
//

type MigrationRecord struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    AppliedAt int64  `json:"applied_at"`
    Checksum  string `json:"checksum"`
    Direction string `json:"direction"` // up / down
}

func (s *Store) SaveMigration(m MigrationRecord) error {
    if m.AppliedAt == 0 {
        m.AppliedAt = time.Now().Unix()
    }
    _, err := s.db.Exec(`
        INSERT INTO migrations (name, applied_at, checksum, direction)
        VALUES (?, ?, ?, ?)
    `, m.Name, m.AppliedAt, m.Checksum, m.Direction)
    return err
}

func (s *Store) LoadMigrations() []MigrationRecord {
    rows, err := s.db.Query(`SELECT id, name, applied_at, checksum, direction FROM migrations ORDER BY applied_at ASC`)
    if err != nil {
        return []MigrationRecord{}
    }
    defer rows.Close()

    var out []MigrationRecord
    for rows.Next() {
        var m MigrationRecord
        if err := rows.Scan(&m.ID, &m.Name, &m.AppliedAt, &m.Checksum, &m.Direction); err != nil {
            continue
        }
        out = append(out, m)
    }

    return out
}

func (s *Store) LoadLastMigration() *MigrationRecord {
    rows, err := s.db.Query(`SELECT id, name, applied_at, checksum, direction FROM migrations ORDER BY applied_at DESC LIMIT 1`)
    if err != nil {
        return nil
    }
    defer rows.Close()

    if rows.Next() {
        var m MigrationRecord
        if err := rows.Scan(&m.ID, &m.Name, &m.AppliedAt, &m.Checksum, &m.Direction); err != nil {
            return nil
        }
        return &m
    }

    return nil
}

func (s *Store) ClearMigrations() error {
    _, err := s.db.Exec(`DELETE FROM migrations`)
    return err
}

//
// -----------------------------------------------------------------------------
// Schema Snapshots Store
// -----------------------------------------------------------------------------
//

type SchemaSnapshot struct {
    ID        int64           `json:"id"`
    Name      string          `json:"name"`
    CreatedAt int64           `json:"created_at"`
    Data      json.RawMessage `json:"data"`
}

func (s *Store) SaveSchemaSnapshot(name string, data interface{}) error {
    raw, _ := json.Marshal(data)
    _, err := s.db.Exec(`
        INSERT INTO schema_snapshots (name, created_at, data)
        VALUES (?, ?, ?)
    `, name, time.Now().Unix(), string(raw))
    return err
}

func (s *Store) LoadSchemaSnapshots() []SchemaSnapshot {
    rows, err := s.db.Query(`SELECT id, name, created_at, data FROM schema_snapshots ORDER BY created_at DESC`)
    if err != nil {
        return []SchemaSnapshot{}
    }
    defer rows.Close()

    var out []SchemaSnapshot
    for rows.Next() {
        var sshot SchemaSnapshot
        var raw string
        if err := rows.Scan(&sshot.ID, &sshot.Name, &sshot.CreatedAt, &raw); err != nil {
            continue
        }
        sshot.Data = json.RawMessage(raw)
        out = append(out, sshot)
    }

    return out
}

func (s *Store) LoadSchemaSnapshotByName(name string) *SchemaSnapshot {
    rows, err := s.db.Query(`SELECT id, name, created_at, data FROM schema_snapshots WHERE name = ? LIMIT 1`, name)
    if err != nil {
        return nil
    }
    defer rows.Close()

    if rows.Next() {
        var sshot SchemaSnapshot
        var raw string
        if err := rows.Scan(&sshot.ID, &sshot.Name, &sshot.CreatedAt, &raw); err != nil {
            return nil
        }
        sshot.Data = json.RawMessage(raw)
        return &sshot
    }

    return nil
}

func (s *Store) DeleteSchemaSnapshot(id int64) error {
    _, err := s.db.Exec(`DELETE FROM schema_snapshots WHERE id = ?`, id)
    return err
}

//
// -----------------------------------------------------------------------------
// Local Cache Store
// -----------------------------------------------------------------------------
//

type LocalCacheEntry struct {
    Key       string `json:"key"`
    Value     string `json:"value"`
    UpdatedAt int64  `json:"updated_at"`
}

func (s *Store) SetCache(key, value string) error {
    now := time.Now().Unix()
    _, err := s.db.Exec(`
        INSERT INTO local_cache (key, value, updated_at)
        VALUES (?, ?, ?)
        ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
    `, key, value, now)
    return err
}

func (s *Store) GetCache(key string) (string, bool) {
    rows, err := s.db.Query(`SELECT value FROM local_cache WHERE key = ? LIMIT 1`, key)
    if err != nil {
        return "", false
    }
    defer rows.Close()

    if rows.Next() {
        var v string
        if err := rows.Scan(&v); err != nil {
            return "", false
        }
        return v, true
    }

    return "", false
}

func (s *Store) DeleteCache(key string) error {
    _, err := s.db.Exec(`DELETE FROM local_cache WHERE key = ?`, key)
    return err
}

func (s *Store) ClearCache() error {
    _, err := s.db.Exec(`DELETE FROM local_cache`)
    return err
}

//
// -----------------------------------------------------------------------------
// Settings Store
// -----------------------------------------------------------------------------
//

func (s *Store) SetSetting(key, value string) error {
    _, err := s.db.Exec(`
        INSERT INTO settings (key, value)
        VALUES (?, ?)
        ON CONFLICT(key) DO UPDATE SET value = excluded.value
    `, key, value)
    return err
}

func (s *Store) GetSetting(key string) (string, bool) {
    rows, err := s.db.Query(`SELECT value FROM settings WHERE key = ? LIMIT 1`, key)
    if err != nil {
        return "", false
    }
    defer rows.Close()

    if rows.Next() {
        var v string
        if err := rows.Scan(&v); err != nil {
            return "", false
        }
        return v, true
    }

    return "", false
}

func (s *Store) DeleteSetting(key string) error {
    _, err := s.db.Exec(`DELETE FROM settings WHERE key = ?`, key)
    return err
}

func (s *Store) LoadAllSettings() map[string]string {
    rows, err := s.db.Query(`SELECT key, value FROM settings`)
    if err != nil {
        return map[string]string{}
    }
    defer rows.Close()

    out := make(map[string]string)
    for rows.Next() {
        var k, v string
        if err := rows.Scan(&k, &v); err != nil {
            continue
        }
        out[k] = v
    }

    return out
}

//
// -----------------------------------------------------------------------------
// File Helpers
// -----------------------------------------------------------------------------
//

func EnsureDir(path string) error {
    return os.MkdirAll(path, 0755)
}

func FileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func RemoveFile(path string) error {
    if !FileExists(path) {
        return nil
    }
    return os.Remove(path)
}

func ProjectSQLitePath(baseDir string, p Project) string {
    if p.SQLitePath != "" {
        return p.SQLitePath
    }
    return filepath.Join(baseDir, "projects", p.Name, "db.sqlite")
}

func EnsureProjectDir(baseDir string, p Project) error {
    dir := filepath.Join(baseDir, "projects", p.Name)
    return EnsureDir(dir)
}

//
// -----------------------------------------------------------------------------
// JSON Helpers
// -----------------------------------------------------------------------------
//

func MustJSON(v interface{}) string {
    b, _ := json.Marshal(v)
    return string(b)
}

func ParseJSON(raw string, v interface{}) error {
    return json.Unmarshal([]byte(raw), v)
}
