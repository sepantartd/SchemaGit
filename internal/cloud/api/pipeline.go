package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func PipelineListHandler(w http.ResponseWriter, req *http.Request) {
    prs, err := store.ListPRs()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(prs)
}

func PipelineDetailHandler(w http.ResponseWriter, req *http.Request) {
    repo := req.URL.Query().Get("repo")
    pr := req.URL.Query().Get("pr")

    plan, _ := store.LoadPlan(repo, atoi(pr))
    job, _ := store.LoadJob(repo + "#" + pr)
    logs, _ := store.LoadLogs(job.ID)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "plan": plan,
        "job":  job,
        "logs": logs,
    })
}
