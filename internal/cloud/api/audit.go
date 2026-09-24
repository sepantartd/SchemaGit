package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func AuditListHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    orgs, err := store.UserOrgs(email)
    if err != nil {
        http.Error(w, "unable to load organizations", http.StatusInternalServerError)
        return
    }
    if len(orgs) == 0 {
        http.Error(w, "no org found", http.StatusBadRequest)
        return
    }

    logs, err := store.ListAudit(orgs[0].ID)
    if err != nil {
        http.Error(w, "unable to load audit logs", http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(logs)
}
