package cloud

import (
    "net/http"
    "github.com/gorilla/mux"
    "schemagit/internal/cloud/api"
)

func NewServer() *mux.Router {
    router := mux.NewRouter()

    // API routes
    router.HandleFunc("/api/pipeline/list", api.PipelineListHandler).Methods("GET")
    router.HandleFunc("/api/pipeline/detail", api.PipelineDetailHandler).Methods("GET")
    router.HandleFunc("/api/pipeline/diff", api.PipelineDiffHandler).Methods("GET")

    // UI routes
    router.HandleFunc("/pipeline", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/pipeline.html")
    })

    router.HandleFunc("/pipeline/pr", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "ui/cloud/pipeline_pr.html")
    })

    router.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))),
    )

    return router
}
