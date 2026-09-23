package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/cloud/store"
)

func OrgCreateHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    var body struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Name == "" {
        http.Error(w, "organization name is required", 400)
        return
    }

    id, err := store.CreateOrg(email, body.Name)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func OrgAddMemberHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    var body struct {
        OrgID string `json:"org_id"`
        Email string `json:"email"`
        Role  string `json:"role"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request", 400)
        return
    }

    if !store.CanManageOrg(email, body.OrgID) {
        http.Error(w, "forbidden", 403)
        return
    }

    if err := store.AddMember(body.OrgID, body.Email, body.Role); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func OrgMembersHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    orgID := req.URL.Query().Get("org")
    if !store.CanAccessOrg(email, orgID) {
        http.Error(w, "forbidden", 403)
        return
    }

    members, err := store.ListMembers(orgID)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(members)
}

func UserOrgsHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    orgs, err := store.UserOrgs(email)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(orgs)
}
