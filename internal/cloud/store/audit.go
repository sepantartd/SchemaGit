package store

import (
    "time"
)

type AuditLog struct {
    ID        int64
    OrgID     string
    Event     string
    Detail    string
    CreatedAt time.Time
}

func AddAudit(orgID, event, detail string) {
    _, _ = DB.Exec(`
        INSERT INTO audit_logs (org_id, event, detail, created_at)
        VALUES (?, ?, ?, datetime('now'))
    `, orgID, event, detail)
}

func ListAudit(orgID string) ([]AuditLog, error) {
    rows, err := DB.Query(`
        SELECT id, org_id, event, detail, created_at
        FROM audit_logs
        WHERE org_id = ?
        ORDER BY created_at DESC
        LIMIT 200
    `, orgID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []AuditLog
    for rows.Next() {
        var a AuditLog
        if err := rows.Scan(&a.ID, &a.OrgID, &a.Event, &a.Detail, &a.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, a)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return out, nil
}
