package cloud

import (
    "encoding/json"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/config"
    "github.com/sepanta/schemagit/internal/log"
)

type Response struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

func StartAPIServer() {
    cfg := config.Default()

    // API endpoints (با Auth)
    http.HandleFunc("/api/diff", withAuth(handleDiff))
    http.HandleFunc("/api/apply", withAuth(handleApply))

    // Webhook GitHub (با Auth)
    http.HandleFunc("/api/webhook/github", withAuth(HandleGitHubWebhook))

    // داشبورد Cloud روی پورت جدا
    go StartDashboardServer(cfg.CloudAPIKey)

    log.Info("Cloud API running on :9090")
    http.ListenAndServe(":9090", nil)
}

func withAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cfg := config.Default()
        key := r.Header.Get("X-API-Key")

        if key == "" || key != cfg.CloudAPIKey {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(Response{
                Ok:    false,
                Error: "unauthorized",
            })
            return
        }

        next(w, r)
    }
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
