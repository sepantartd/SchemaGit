package cloud

import (
    "html/template"
    "net/http"
    "strconv"

    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"
)

type DashboardData struct {
    Title    string
    APIKey   string
    Projects []store.Project
}

type ProjectDetailData struct {
    Title    string
    Project  *store.Project
    Logs     []store.CloudLog
}

var dashboardTpl = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))
var projectTpl   = template.Must(template.ParseFiles("ui/cloud/project.html"))

func StartDashboardServer(apiKey string) {

    // صفحهٔ اصلی داشبورد
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := DashboardData{
            Title:    "SchemaGit Cloud Dashboard",
            APIKey:   apiKey,
            Projects: GetProjects(),
        }
        dashboardTpl.Execute(w, data)
    })

    // صفحهٔ جزئیات پروژه
    http.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
        pidStr := r.URL.Query().Get("id")
        if pidStr == "" {
            http.Error(w, "missing project id", 400)
            return
        }

        pid, err := strconv.ParseInt(pidStr, 10, 64)
        if err != nil {
            http.Error(w, "invalid project id", 400)
            return
        }

        project := GetProjectByID(pid)
        if project == nil {
            http.Error(w, "project not found", 404)
            return
        }

        logs := GetLogsByProject(pid)

        data := ProjectDetailData{
            Title:   "Project Detail",
            Project: project,
            Logs:    logs,
        }

        projectTpl.Execute(w, data)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
