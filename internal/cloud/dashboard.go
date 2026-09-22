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
    Title   string
    Project *store.Project
    Logs    []store.CloudLog
}

var dashboardTpl        = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))
var projectTpl          = template.Must(template.ParseFiles("ui/cloud/project.html"))
var projectSettingsTpl  = template.Must(template.ParseFiles("ui/cloud/project_settings.html"))

func StartDashboardServer(apiKey string) {

    //
    // Dashboard (list of projects)
    //
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := DashboardData{
            Title:    "SchemaGit Cloud Dashboard",
            APIKey:   apiKey,
            Projects: GetProjects(),
        }
        dashboardTpl.Execute(w, data)
    })

    //
    // Project Detail Page
    //
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

    //
    // Project Settings Page
    //
    http.HandleFunc("/project/settings", func(w http.ResponseWriter, r *http.Request) {

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

        // Save form
        if r.Method == "POST" {

            r.ParseForm()

            project.DBType = r.Form.Get("db_type")
            project.SQLitePath = r.Form.Get("sqlite_path")

            project.PGHost = r.Form.Get("pg_host")
            project.PGPort, _ = strconv.Atoi(r.Form.Get("pg_port"))
            project.PGUser = r.Form.Get("pg_user")
            project.PGPass = r.Form.Get("pg_pass")
            project.PGName = r.Form.Get("pg_name")

            project.MYHost = r.Form.Get("my_host")
            project.MYPort, _ = strconv.Atoi(r.Form.Get("my_port"))
            project.MYUser = r.Form.Get("my_user")
            project.MYPass = r.Form.Get("my_pass")
            project.MYName = r.Form.Get("my_name")

            UpdateProject(project)

            http.Redirect(w, r, "/project?id="+pidStr, 302)
            return
        }

        // Show form
        data := struct {
            Title   string
            Project *store.Project
        }{
            Title:   "Project Settings",
            Project: project,
        }

        projectSettingsTpl.Execute(w, data)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
