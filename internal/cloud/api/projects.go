package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func CreateProjectHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    if !store.CanCreateProject(email) {
        http.Error(w, "project limit reached for your plan", http.StatusForbidden)
        return
    }

    var body struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Name == "" {
        http.Error(w, "project name is required", http.StatusBadRequest)
        return
    }

    if err := store.CreateProject(email, body.Name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
