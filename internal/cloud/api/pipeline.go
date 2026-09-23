package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "schemagit/internal/cloud/store"
)

func atoi(s string) int {
    n, _ := strconv.Atoi(s)
    return n
}

type PRItem struct {
    Repo   string `json:"repo"`
    Number int    `json:"number"`
    Status string `json:"status"`
}

// List PRs (already used by /pipeline)
func PipelineListHandler(w http.ResponseWriter, req *http.Request) {
    prs, err := store.ListPRs()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(prs)
}

// Detail: plan + job + logs (already used by /pipeline/pr)
func PipelineDetailHandler(w http.ResponseWriter, req *http.Request) {
    repo := req.URL.Query().Get("repo")
    pr := req.URL.Query().Get("pr")

    plan, _ := store.LoadPlan(repo, atoi(pr))
    job, _ := store.LoadJob(repo + "#" + pr)
    logs, _ := store.LoadLogs(job.ID)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "plan": plan,
        "job":  job,
        // logs: array of log lines used by the enhanced Job Logs UI
        "logs": logs,
    })
}

// NEW: Diff endpoint for Diff Viewer
func PipelineDiffHandler(w http.ResponseWriter, req *http.Request) {
    repo := req.URL.Query().Get("repo")
    pr := req.URL.Query().Get("pr")

    plan, err := store.LoadPlan(repo, atoi(pr))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // We assume the MigrationPlan contains a Diff field.
    // If not, you can adjust this to return the whole plan.
    json.NewEncoder(w).Encode(map[string]interface{}{
        "plan": plan,
    })
}

func PipelinePlanHandler(w http.ResponseWriter, req *http.Request) {
    repo := req.URL.Query().Get("repo")
    pr := req.URL.Query().Get("pr")

    plan, err := store.LoadPlan(repo, atoi(pr))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(plan)
}
