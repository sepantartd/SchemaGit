package store

import "testing"

func TestOrgMembershipAndAccess(t *testing.T) {
    db, err := InitTempDB()
    if err != nil {
        t.Fatal(err)
    }
    DB = db

    if err := CreateUser("owner@example.com", "secret"); err != nil {
        t.Fatalf("create user: %v", err)
    }
    if err := CreateUser("member@example.com", "secret"); err != nil {
        t.Fatalf("create user: %v", err)
    }

    orgID, err := CreateOrg("owner@example.com", "Acme")
    if err != nil {
        t.Fatalf("create org: %v", err)
    }
    if err := AddMember(orgID, "member@example.com", "member"); err != nil {
        t.Fatalf("add member: %v", err)
    }

    if !CanAccessOrg("owner@example.com", orgID) {
        t.Fatal("owner should access org")
    }
    if !CanAccessOrg("member@example.com", orgID) {
        t.Fatal("member should access org")
    }
    if CanAccessOrg("other@example.com", orgID) {
        t.Fatal("non-member should not access org")
    }
    if !CanManageOrg("owner@example.com", orgID) {
        t.Fatal("owner should manage org")
    }
    if CanManageOrg("member@example.com", orgID) {
        t.Fatal("member should not manage org")
    }
}

func InitTempDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return nil, err
    }
    DB = db
    if err := initTables(); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}
