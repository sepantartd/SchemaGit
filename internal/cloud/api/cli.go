package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func CLISyncHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("X-API-Token")
    email, err := store.ValidateAPIToken(token)
    if err != nil {
        http.Error(w, "invalid API token", 401)
        return
    }

    projects := store.UserProjects(email)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "email":    email,
        "projects": projects,
    })
}

func CLISyncProjectHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("X-API-Token")
    email, err := store.ValidateAPIToken(token)
    if err != nil {
        http.Error(w, "invalid API token", 401)
        return
    }

    project := req.URL.Query().Get("project")

    schema, _ := store.LoadProjectSchema(email, project)
    migrations, _ := store.LoadProjectMigrations(email, project)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "schema":     schema,
        "migrations": migrations,
    })
}
