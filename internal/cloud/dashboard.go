package cloud

import (
    "html/template"
    "net/http"
    "strconv"

    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"
)

type DashboardData struct {
    Title          string
    APIKey         string
    Logs           []store.CloudLog
    Projects       []store.Project
    CurrentProject int64
}

var dashboardTpl = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))

func StartDashboardServer(apiKey string) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        pidStr := r.URL.Query().Get("project")
        var pid int64 = 0

        if pidStr != "" {
            pid, _ = strconv.ParseInt(pidStr, 10, 64)
        }

        logs := GetLogsByProject(pid)
        projects := GetProjects()

        data := DashboardData{
            Title:          "SchemaGit Cloud Dashboard",
            APIKey:         apiKey,
            Logs:           logs,
            Projects:       projects,
            CurrentProject: pid,
        }

        dashboardTpl.Execute(w, data)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
