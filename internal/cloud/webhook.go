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
}    }

    mac := hmac.New(sha1.New, []byte(secret))
    mac.Write(body)
    expected := "sha1=" + hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(expected), []byte(signature))
}

//
// Find project by webhook ID
//

func findProjectByWebhook(webhookID string) int64 {
    projects := GetProjects()
    for _, p := range projects {
        if p.GitHubWebhookURL != "" && strings.Contains(p.GitHubWebhookURL, webhookID) {
            return p.ID
        }
    }
    return 0
}

//
// Create Job in Cloud
//

func CreateJob(apiKey, jobType string, payload map[string]string) {
    data := map[string]interface{}{
        "type":    jobType,
        "payload": payload,
        "created": time.Now().Unix(),
    }

    apiPost("/jobs/create", apiKey, data)
}

//
// Webhook Server
//

func StartWebhookServer(apiKey string) {
    log.Info("Webhook server running on :8090")

    http.HandleFunc("/webhook/github", func(w http.ResponseWriter, r *http.Request) {

        // Read body
        rawBody, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "cannot read body", 400)
            return
        }

        // Extract signature
        signature := r.Header.Get("X-Hub-Signature")
        if signature == "" {
            http.Error(w, "missing signature", 401)
            return
        }

        // Extract webhook ID
        webhookID := r.URL.Query().Get("id")
        if webhookID == "" {
            http.Error(w, "missing webhook id", 400)
            return
        }

        // Find project
        projectID := findProjectByWebhook(webhookID)
        if projectID == 0 {
            http.Error(w, "project not found", 404)
            return
        }

        project := GetProjectByID(projectID)
        if project == nil {
            http.Error(w, "project not found", 404)
            return
        }

        // Validate signature
        if !validateGitHubSignature(project.GitHubSecret, rawBody, signature) {
            http.Error(w, "invalid signature", 401)
            return
        }

        // Parse payload
        var payload GitHubPushPayload
        json.Unmarshal(rawBody, &payload)

        repo := payload.Repository.CloneURL
        branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

        // Create job
        CreateJob(apiKey, "github_push", map[string]string{
            "repo":       repo,
            "branch":     branch,
            "project_id": fmtInt(projectID),
        })

        log.Info("Webhook received → Job created for project " + project.Name)

        w.WriteHeader(200)
        w.Write([]byte(`{"ok":true}`))
    })

    http.ListenAndServe(":8090", nil)
}

//
// API Helpers
//

func apiPost(path, apiKey string, data map[string]interface{}) []byte {
    body, _ := json.Marshal(data)

    req, _ := http.NewRequest("POST", "http://localhost:9090"+path, strings.NewReader(string(body)))
    req.Header.Set("X-API-Key", apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()

    out, _ := io.ReadAll(resp.Body)
    return out
}

func apiCall(path, apiKey string) *http.Response {
    req, _ := http.NewRequest("GET", "http://localhost:9090"+path, nil)
    req.Header.Set("X-API-Key", apiKey)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil
    }
    return resp
}

func fmtInt(v int64) string {
    return strconv.FormatInt(v, 10)
}
