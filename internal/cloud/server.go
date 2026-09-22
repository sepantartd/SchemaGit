package cloud

import (
    "encoding/json"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/config"
    "github.com/sepanta/schemagit/internal/log"
)

type Response struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

func StartAPIServer() {
    cfg := config.Default()

    // init logs DB
    InitLogs(".cloud_logs.db")

    // API endpoints
    http.HandleFunc("/api/diff", withAuth(handleDiff))
    http.HandleFunc("/api/apply", withAuth(handleApply))
    http.HandleFunc("/api/webhook/github", withAuth(HandleGitHubWebhook))

    // dashboard
    go StartDashboardServer(cfg.CloudAPIKey)

    log.Info("Cloud API running on :9090")
    http.ListenAndServe(":9090", nil)
}

func withAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cfg := config.Default()
        key := r.Header.Get("X-API-Key")

        if key == "" || key != cfg.CloudAPIKey {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(Response{
                Ok:    false,
                Error: "unauthorized",
            })
            return
        }

        next(w, r)
    }
}

func handleDiff(w http.ResponseWriter, r *http.Request) {
    err := cli.RunDiff()
    AddLog("diff", err == nil, errString(err))
    respond(w, err)
}

func handleApply(w http.ResponseWriter, r *http.Request) {
    err := cli.RunApply()
    AddLog("apply", err == nil, errString(err))
    respond(w, err)
}

func HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
    err1 := cli.RunDiff()
    err2 := cli.RunApply()

    success := err1 == nil && err2 == nil
    errMsg := errString(err1) + " | " + errString(err2)

    AddLog("webhook", success, errMsg)

    respond(w, nil)
}

func respond(w http.ResponseWriter, err error) {
    resp := Response{Ok: err == nil}
    if err != nil {
        resp.Error = err.Error()
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func errString(err error) string {
    if err == nil {
        return ""
    }
    return err.Error()
}
