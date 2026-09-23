package cloud

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "strings"
    "time"

    "github.com/sepanta/schemagit/internal/log"
)

func StartWebhookServer(apiKey string) {
    http.HandleFunc("/webhook/github", func(w http.ResponseWriter, r *http.Request) {

        rawBody, _ := io.ReadAll(r.Body)
        signature := r.Header.Get("X-Hub-Signature")

        webhookID := r.URL.Query().Get("id")
        projectID := findProjectByWebhook(webhookID)

        project := GetProjectByID(projectID)

        if !validateGitHubSignature(project.GitHubSecret, rawBody, signature) {
            http.Error(w, "invalid signature", 401)
            return
        }

        var payload struct {
            Ref        string `json:"ref"`
            Repository struct {
                CloneURL string `json:"clone_url"`
            } `json:"repository"`
        }

        json.Unmarshal(rawBody, &payload)

        branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

        CreateJob(apiKey, "github_push", map[string]string{
            "repo":       payload.Repository.CloneURL,
            "branch":     branch,
            "project_id": fmtInt(projectID),
        })

        w.Write([]byte(`{"ok":true}`))
    })

    http.ListenAndServe(":8090", nil)
}
