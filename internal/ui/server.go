package ui

import (
    "html/template"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/log"
)

func StartServer() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/diff", handleDiff)
    http.HandleFunc("/apply", handleApply)

    log.Info("UI server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("ui/templates/home.html"))
    t.Execute(w, nil)
}

func handleDiff(w http.ResponseWriter, r *http.Request) {
    err := cli.RunDiff()
    t := template.Must(template.ParseFiles("ui/templates/diff.html"))

    data := struct {
        Error string
    }{}

    if err != nil {
        data.Error = err.Error()
    }

    t.Execute(w, data)
}

func handleApply(w http.ResponseWriter, r *http.Request) {
    err := cli.RunApply()
    t := template.Must(template.ParseFiles("ui/templates/apply.html"))

    data := struct {
        Error string
    }{}

    if err != nil {
        data.Error = err.Error()
    }

    t.Execute(w, data)
}
