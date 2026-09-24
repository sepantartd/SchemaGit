package cloud

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/config"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"

    "schemagit/internal/cloud/api"
)

type Response struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

var logsStore *store.Store
var projStore *store.Store

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

func GetLogsByProject(pid int64) []store.CloudLog {
    if logsStore == nil {
        return []store.CloudLog{}
    }
    return logsStore.LoadLogsByProject(pid)
}

func GetProjects() []store.Project {
    if projStore == nil {
        return []store.Project{}
    }
    return projStore.LoadProjects()
}

func StartAPIServer() {
    cfg := config.Default()

    InitLogs(".cloud_logs.db")
    InitProjects(".cloud_projects.db")

    http.HandleFunc("/api/pipeline/plan", api.PipelinePlanHandler)
    http.HandleFunc("/api/project/overview", api.ProjectOverviewHandler)
    http.HandleFunc("/api/tokens/create", api.TokenCreateHandler)
    http.HandleFunc("/api/tokens/list", api.TokenListHandler)
    http.HandleFunc("/api/tokens/delete", api.TokenDeleteHandler)
    http.HandleFunc("/api/org/create", api.OrgCreateHandler)
    http.HandleFunc("/api/org/add_member", api.OrgAddMemberHandler)
    http.HandleFunc("/api/org/members", api.OrgMembersHandler)
    http.HandleFunc("/api/org/list", api.UserOrgsHandler)
    http.HandleFunc("/api/webhooks/add", api.WebhookAddHandler)
    http.HandleFunc("/api/webhooks/list", api.WebhookListHandler)
    http.HandleFunc("/api/notifications/set", api.NotificationSetHandler)
    http.HandleFunc("/api/notifications/get", api.NotificationGetHandler)
    http.HandleFunc("/api/audit/list", api.AuditListHandler)
    http.HandleFunc("/api/cli/sync", api.CLISyncHandler)
    http.HandleFunc("/api/cli/sync/project", api.CLISyncProjectHandler)
    http.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/tokens.html")
    })
    http.HandleFunc("/org", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/org.html")
    })
    http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/webhooks.html")
    })
    http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/notifications.html")
    })
    http.HandleFunc("/audit", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/audit.html")
    })
    http.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/project_overview.html")
    })

    http.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path

        if strings.HasSuffix(path, "/diff") {
            withAuth(withProjectFromURL(handleDiff))(w, r)
            return
        }

        if strings.HasSuffix(path, "/apply") {
            withAuth(withProjectFromURL(handleApply))(w, r)
            return
        }

        if strings.HasSuffix(path, "/webhook/github") {
            withAuth(withProjectFromURL(HandleGitHubWebhook))(w, r)
            return
        }

        respond(w, errors.New("unknown project endpoint"))
    })

    http.HandleFunc("/agent/logs", RequireAPIKey(AgentLogsAPI))
    http.HandleFunc("/dashboard/agent/logs", ProjectAgentLogs)

    go StartDashboardServer(cfg.CloudAPIKey)

    handler := http.DefaultServeMux
    handler = LoggingMiddleware(handler)
    handler = SecurityHeadersMiddleware(handler)
    handler = CORSMiddleware(handler)
    handler = RecoveryMiddleware(handler)
    handler = RateLimitMiddleware(100, time.Minute)(handler)

    log.Info("Cloud API running on :9090")
    http.ListenAndServe(":9090", handler)
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

func withProjectFromURL(next func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) < 4 {
            respond(w, errors.New("invalid project path"))
            return
        }

        pid, err := strconv.ParseInt(parts[3], 10, 64)
        if err != nil {
            respond(w, errors.New("invalid project id"))
            return
        }

        next(w, r, pid)
    }
}

func handleDiff(w http.ResponseWriter, r *http.Request, projectID int64) {
    project := GetProjectByID(projectID)
    if project == nil {
        respond(w, errors.New("project not found"))
        return
    }

    config.SetDBForProject(project)

    err := cli.RunDiff()
    AddLog("diff", err == nil, errString(err), projectID)
    respond(w, err)
}

func handleApply(w http.ResponseWriter, r *http.Request, projectID int64) {
    project := GetProjectByID(projectID)
    if project == nil {
        respond(w, errors.New("project not found"))
        return
    }

    config.SetDBForProject(project)

    err := cli.RunApply()
    AddLog("apply", err == nil, errString(err), projectID)
    respond(w, err)
}

func HandleGitHubWebhook(w http.ResponseWriter, r *http.Request, projectID int64) {
    project := GetProjectByID(projectID)
    if project == nil {
        respond(w, errors.New("project not found"))
        return
    }

    config.SetDBForProject(project)

    err1 := cli.RunDiff()
    err2 := cli.RunApply()

    success := err1 == nil && err2 == nil
    errMsg := errString(err1) + " | " + errString(err2)

    AddLog("webhook", success, errMsg, projectID)

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
