package api

import (
    "net/http"

    "github.com/gorilla/mux"
    "schemagit/internal/cloud/runner"
)

func NewRouter(r *runner.Runner) *mux.Router {
    router := mux.NewRouter()

    // ثبت API اجرای Migration
    router.HandleFunc(
        "/projects/{id}/migrations/run",
        RunMigrationHandler(r),
    ).Methods("POST")

    return router
}
