package cloud

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func render(w http.ResponseWriter, name string, data interface{}) {
    path := filepath.Join("dashboard", name)

    tmpl, err := template.ParseFiles(path)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
