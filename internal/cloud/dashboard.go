package cloud

import (
    "html/template"
    "net/http"

    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"
)

type DashboardData struct {
    Title    string
    APIKey   string
    Logs     []store.CloudLog
    Projects []store.Project
}

var dashboardTpl = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))

func StartDashboardServer(apiKey string) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logs := GetLogs()
        projects := GetProjects()

        data := DashboardData{
            Title:    "SchemaGit Cloud Dashboard",
            APIKey:   apiKey,
            Logs:     logs,
            Projects: projects,
        }

        dashboardTpl.Execute(w, data)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
