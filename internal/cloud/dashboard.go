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
    Members []store.ProjectMember
}

var dashboardTpl       = template.Must(template.ParseFiles("ui/cloud/dashboard.html"))
var projectTpl         = template.Must(template.ParseFiles("ui/cloud/project.html"))
var projectSettingsTpl = template.Must(template.ParseFiles("ui/cloud/project_settings.html"))
var projectCreateTpl   = template.Must(template.ParseFiles("ui/cloud/project_create.html"))
var projectMembersTpl  = template.Must(template.ParseFiles("ui/cloud/project_members.html"))

func StartDashboardServer(apiKey string) {

    // Dashboard (list of projects)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := DashboardData{
            Title:    "SchemaGit Cloud Dashboard",
            APIKey:   apiKey,
            Projects: GetProjects(),
        }
        dashboardTpl.Execute(w, data)
    })

    // Project Detail Page
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
        members := GetMembersByProject(pid)

        data := ProjectDetailData{
            Title:   "Project Detail",
            Project: project,
            Logs:    logs,
            Members: members,
        }

        projectTpl.Execute(w, data)
    })

    // Project Settings Page
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

            project.GitHubWebhookURL = r.Form.Get("github_webhook_url")
            project.GitHubSecret = r.Form.Get("github_secret")

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

    // Project Create Page
    http.HandleFunc("/project/create", func(w http.ResponseWriter, r *http.Request) {

        if r.Method == "POST" {

            r.ParseForm()

            name := r.Form.Get("name")
            if name == "" {
                http.Error(w, "name is required", 400)
                return
            }

            p := store.Project{
                Name:       name,
                DBType:     r.Form.Get("db_type"),
                SQLitePath: r.Form.Get("sqlite_path"),

                PGHost: r.Form.Get("pg_host"),
                PGPort: func() int {
                    v, _ := strconv.Atoi(r.Form.Get("pg_port"))
                    return v
                }(),
                PGUser: r.Form.Get("pg_user"),
                PGPass: r.Form.Get("pg_pass"),
                PGName: r.Form.Get("pg_name"),

                MYHost: r.Form.Get("my_host"),
                MYPort: func() int {
                    v, _ := strconv.Atoi(r.Form.Get("my_port"))
                    return v
                }(),
                MYUser: r.Form.Get("my_user"),
                MYPass: r.Form.Get("my_pass"),
                MYName: r.Form.Get("my_name"),

                GitHubWebhookURL: r.Form.Get("github_webhook_url"),
                GitHubSecret:     r.Form.Get("github_secret"),
            }

            projStore.SaveProject(p)

            http.Redirect(w, r, "/", 302)
            return
        }

        data := struct {
            Title string
        }{
            Title: "Create Project",
        }

        projectCreateTpl.Execute(w, data)
    })

    // Project Delete
    http.HandleFunc("/project/delete", func(w http.ResponseWriter, r *http.Request) {
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

        DeleteProject(pid)

        http.Redirect(w, r, "/", 302)
    })

    // Project Members Page
    http.HandleFunc("/project/members", func(w http.ResponseWriter, r *http.Request) {
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

        members := GetMembersByProject(pid)

        data := struct {
            Title   string
            Project *store.Project
            Members []store.ProjectMember
        }{
            Title:   "Project Members",
            Project: project,
            Members: members,
        }

        projectMembersTpl.Execute(w, data)
    })

    // Add Member
    http.HandleFunc("/project/members/add", func(w http.ResponseWriter, r *http.Request) {
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

        if r.Method != "POST" {
            http.Error(w, "method not allowed", 405)
            return
        }

        r.ParseForm()
        email := r.Form.Get("email")
        role := r.Form.Get("role")

        if email == "" || role == "" {
            http.Error(w, "email and role required", 400)
            return
        }

        _ = AddMemberToProject(pid, email, role)

        http.Redirect(w, r, "/project/members?id="+pidStr, 302)
    })

    // Delete Member
    http.HandleFunc("/project/members/delete", func(w http.ResponseWriter, r *http.Request) {
        midStr := r.URL.Query().Get("member_id")
        pidStr := r.URL.Query().Get("project_id")

        if midStr == "" || pidStr == "" {
            http.Error(w, "missing member_id or project_id", 400)
            return
        }

        mid, err := strconv.ParseInt(midStr, 10, 64)
        if err != nil {
            http.Error(w, "invalid member id", 400)
            return
        }

        _ = DeleteMember(mid)

        http.Redirect(w, r, "/project/members?id="+pidStr, 302)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}            APIKey:   apiKey,
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

    //
    // Project Create Page
    //
    http.HandleFunc("/project/create", func(w http.ResponseWriter, r *http.Request) {

        if r.Method == "POST" {

            r.ParseForm()

            name := r.Form.Get("name")
            if name == "" {
                http.Error(w, "name is required", 400)
                return
            }

            p := store.Project{
                Name:       name,
                DBType:     r.Form.Get("db_type"),
                SQLitePath: r.Form.Get("sqlite_path"),

                PGHost: r.Form.Get("pg_host"),
                PGPort: func() int {
                    v, _ := strconv.Atoi(r.Form.Get("pg_port"))
                    return v
                }(),
                PGUser: r.Form.Get("pg_user"),
                PGPass: r.Form.Get("pg_pass"),
                PGName: r.Form.Get("pg_name"),

                MYHost: r.Form.Get("my_host"),
                MYPort: func() int {
                    v, _ := strconv.Atoi(r.Form.Get("my_port"))
                    return v
                }(),
                MYUser: r.Form.Get("my_user"),
                MYPass: r.Form.Get("my_pass"),
                MYName: r.Form.Get("my_name"),
            }

            projStore.SaveProject(p)

            http.Redirect(w, r, "/", 302)
            return
        }

        data := struct {
            Title string
        }{
            Title: "Create Project",
        }

        projectCreateTpl.Execute(w, data)
    })

    //
    // Project Delete
    //
    http.HandleFunc("/project/delete", func(w http.ResponseWriter, r *http.Request) {
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

        DeleteProject(pid)

        http.Redirect(w, r, "/", 302)
    })

    log.Info("Cloud Dashboard running on :8081")
    http.ListenAndServe(":8081", nil)
}
