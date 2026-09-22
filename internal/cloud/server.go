package cloud

import (
    "encoding/json"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/log"
)

type Response struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

func StartAPIServer() {
    http.HandleFunc("/api/diff", handleDiff)
    http.HandleFunc("/api/apply", handleApply)

    log.Info("Cloud API running on :9090")
    http.ListenAndServe(":9090", nil)
}

func handleDiff(w http.ResponseWriter, r *http.Request) {
    err := cli.RunDiff()

    resp := Response{Ok: err == nil}
    if err != nil {
        resp.Error = err.Error()
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func handleApply(w http.ResponseWriter, r *http.Request) {
    err := cli.RunApply()

    resp := Response{Ok: err == nil}
    if err != nil {
        resp.Error = err.Error()
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
