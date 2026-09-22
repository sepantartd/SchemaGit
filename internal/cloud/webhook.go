package cloud

import (
    "encoding/json"
    "net/http"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/log"
)

type WebhookResponse struct {
    Ok    bool   `json:"ok"`
    Error string `json:"error,omitempty"`
}

func HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
    var payload map[string]interface{}

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(WebhookResponse{
            Ok:    false,
            Error: "invalid json",
        })
        return
    }

    log.Info("Webhook received → running diff/apply")

    if err := cli.RunDiff(); err != nil {
        json.NewEncoder(w).Encode(WebhookResponse{
            Ok:    false,
            Error: err.Error(),
        })
        return
    }

    if err := cli.RunApply(); err != nil {
        json.NewEncoder(w).Encode(WebhookResponse{
            Ok:    false,
            Error: err.Error(),
        })
        return
    }

    json.NewEncoder(w).Encode(WebhookResponse{
        Ok: true,
    })
}
