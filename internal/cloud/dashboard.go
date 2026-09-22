package cloud

import (
    "html/template"
    "net/http"

    "github.com/sepanta/schemagit/internal/log"
)

type DashboardData struct {
    Title   string
    APIKey  string
    Logs    []CloudLog
}

var dashboardTpl = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))

func StartDashboardServer(apiKey string) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logs := GetLogs()

        data := DashboardData{
            Title:  "SchemaGit Cloud Dashboard",
            APIKey: apiKey,
            Logs:   logs,
        }

        dashboardTpl.Execute(w, data)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
