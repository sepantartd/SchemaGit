package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func ProjectOverviewHandler(w http.ResponseWriter, req *http.Request) {
    projectID := req.URL.Query().Get("project")

    prs, _ := store.ListPRs()
    jobs, _ := store.ListJobs(projectID)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "prs":  prs,
        "jobs": jobs,
    })
}
