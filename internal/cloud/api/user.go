package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/cloud/store"
)

func UserInfoHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    plan, err := store.GetBillingPlan(email)
    if err != nil {
        http.Error(w, "unable to load billing plan", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "email":      email,
        "plan":       plan.Plan,
        "billing":    map[string]string{"status": "active"},
        "projects":   store.CountProjects(email),
        "migrations": store.CountMigrations(email),
        "limits": map[string]int{
            "projects":   plan.MaxProjects,
            "migrations": plan.MaxMigrations,
        },
    })
}

func ChangePasswordHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    var body struct {
        Old string `json:"old"`
        New string `json:"new"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if err := store.ChangePassword(email, body.Old, body.New); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DeleteAccountHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    store.DeleteUser(email)
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
