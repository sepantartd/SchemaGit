package store

import (
    "errors"

    "github.com/google/uuid"
)

type Org struct {
    ID    string
    Name  string
    Owner string
}

type OrgMember struct {
    OrgID string
    Email string
    Role  string
}

func CreateOrg(owner, name string) (string, error) {
    id := uuid.New().String()

    _, err := DB.Exec(`
        INSERT INTO orgs (id, name, owner)
        VALUES (?, ?, ?)
    `, id, name, owner)
    if err != nil {
        return "", err
    }

    _, _ = DB.Exec(`
        INSERT INTO org_members (org_id, email, role)
        VALUES (?, ?, 'owner')
    `, id, owner)

    return id, nil
}

func AddMember(orgID, email, role string) error {
    if role != "owner" && role != "admin" && role != "member" {
        return errors.New("invalid role")
    }

    _, err := DB.Exec(`
        INSERT INTO org_members (org_id, email, role)
        VALUES (?, ?, ?)
    `, orgID, email, role)

    return err
}

func ListMembers(orgID string) ([]OrgMember, error) {
    rows, err := DB.Query(`
        SELECT org_id, email, role FROM org_members WHERE org_id = ?
    `, orgID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := []OrgMember{}
    for rows.Next() {
        var m OrgMember
        if err := rows.Scan(&m.OrgID, &m.Email, &m.Role); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    return out, nil
}

func UserOrgs(email string) ([]Org, error) {
    rows, err := DB.Query(`
        SELECT o.id, o.name, o.owner
        FROM orgs o
        JOIN org_members m ON o.id = m.org_id
        WHERE m.email = ?
    `, email)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := []Org{}
    for rows.Next() {
        var o Org
        if err := rows.Scan(&o.ID, &o.Name, &o.Owner); err != nil {
            return nil, err
        }
        out = append(out, o)
    }
    return out, nil
}

func CanAccessOrg(email, orgID string) bool {
    row := DB.QueryRow(`
        SELECT 1 FROM org_members WHERE org_id = ? AND email = ?
    `, orgID, email)
    var v int
    if err := row.Scan(&v); err != nil {
        return false
    }
    return true
}

func CanManageOrg(email, orgID string) bool {
    row := DB.QueryRow(`
        SELECT role FROM org_members WHERE org_id = ? AND email = ?
    `, orgID, email)
    var role string
    if err := row.Scan(&role); err != nil {
        return false
    }
    return role == "owner" || role == "admin"
}
