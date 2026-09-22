package cloud

import (
    "encoding/json"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/config"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"
)

type Response struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

var logsStore *store.Store
var projStore *store.Store

//
// Init
//

func InitLogs(path string) {
    st, err := store.Open(path)
    if err != nil {
        log.Error("cannot open logs store: " + err.Error())
        return
    }
    logsStore = st
}

func InitProjects(path string) {
    st, err := store.Open(path)
    if err != nil {
        log.Error("cannot open projects store: " + err.Error())
        return
    }
    projStore = st
}

//
// Helpers
//

func AddLog(t string, success bool, errMsg string, projectID int64) {
    if logsStore == nil {
        return
    }

    entry := store.CloudLog{
        Type:      t,
        Success:   success,
        Error:     errMsg,
        Timestamp: time.Now().Unix(),
        ProjectID: projectID,
    }

    logsStore.SaveCloudLog(entry)
}

func GetLogs() []store.CloudLog {
    if logsStore == nil {
        return []store.CloudLog{}
    }
    return logsStore.LoadCloudLogs()
}

func GetProjects() []store.Project {
    if projStore == nil {
        return []store.Project{}
    }
    return projStore.LoadProjects()
}

//
// API Server
//

func StartAPIServer() {
    cfg := config.Default()

    InitLogs(".cloud_logs.db")
    InitProjects(".cloud_projects.db")

    http.HandleFunc("/api/diff", withAuth(handleDiff))
    http.HandleFunc("/api/apply", withAuth(handleApply))
    http.HandleFunc("/api/webhook/github", withAuth(HandleGitHubWebhook))

    go StartDashboardServer(cfg.CloudAPIKey)

    log.Info("Cloud API running on :9090")
    http.ListenAndServe(":9090", nil)
}

//
// Auth
//

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

//
// Handlers
//

func handleDiff(w http.ResponseWriter, r *http.Request) {
    err := cli.RunDiff()
    AddLog("diff", err == nil, errString(err), 0)
    respond(w, err)
}

func handleApply(w http.ResponseWriter, r *http.Request) {
    err := cli.RunApply()
    AddLog("apply", err == nil, errString(err), 0)
    respond(w, err)
}

func HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
    err1 := cli.RunDiff()
    err2 := cli.RunApply()

    success := err1 == nil && err2 == nil
    errMsg := errString(err1) + " | " + errString(err2)

    AddLog("webhook", success, errMsg, 0)

    respond(w, nil)
}

//
// Response helper
//

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
