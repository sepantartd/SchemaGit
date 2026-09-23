package cloud

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "strings"
)

func validateGitHubSignature(secret string, body []byte, signature string) bool {
    if !strings.HasPrefix(signature, "sha256=") {
        return false
    }

    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(expected), []byte(signature))
}

func StartWebhookServer(apiKey string) {
    http.HandleFunc("/webhook/github", func(w http.ResponseWriter, r *http.Request) {

        rawBody, _ := io.ReadAll(r.Body)
        signature := r.Header.Get("X-Hub-Signature-256")

        webhookID := r.URL.Query().Get("id")
        projectID := findProjectByWebhook(webhookID)
        project := GetProjectByID(projectID)

        if !validateGitHubSignature(project.GitHubSecret, rawBody, signature) {
            http.Error(w, "invalid signature", 401)
            return
        }

        event := r.Header.Get("X-GitHub-Event")

        switch event {

        case "push":
            handlePushWebhook(apiKey, projectID, rawBody)

        case "pull_request":
            handlePRWebhook(apiKey, projectID, rawBody)

        case "pull_request_review":
            handlePRReviewWebhook(apiKey, projectID, rawBody)

        default:
            w.Write([]byte(`{"ignored":true}`))
            return
        }

        w.Write([]byte(`{"ok":true}`))
    })

    http.ListenAndServe(":8090", nil)
}

func handlePushWebhook(apiKey string, projectID int64, raw []byte) {
    var payload struct {
        Ref        string `json:"ref"`
        Repository struct {
            CloneURL string `json:"clone_url"`
        } `json:"repository"`
    }

    json.Unmarshal(raw, &payload)
    branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

    CreateJob(apiKey, "github_push", map[string]string{
        "repo":       payload.Repository.CloneURL,
        "branch":     branch,
        "project_id": fmtInt(projectID),
    })
}

func handlePRWebhook(apiKey string, projectID int64, raw []byte) {
    var payload struct {
        Action     string `json:"action"`
        PullRequest struct {
            Head struct {
                Ref     string `json:"ref"`
                Repo    struct {
                    CloneURL string `json:"clone_url"`
                } `json:"repo"`
            } `json:"head"`
        } `json:"pull_request"`
    }

    json.Unmarshal(raw, &payload)

    if payload.Action != "opened" && payload.Action != "synchronize" {
        return
    }

    CreateJob(apiKey, "github_pr_diff", map[string]string{
        "repo":       payload.PullRequest.Head.Repo.CloneURL,
        "branch":     payload.PullRequest.Head.Ref,
        "project_id": fmtInt(projectID),
    })
}

func handlePRReviewWebhook(apiKey string, projectID int64, raw []byte) {
    var payload struct {
        Action string `json:"action"`
        PullRequest struct {
            Head struct {
                Ref     string `json:"ref"`
                Repo    struct {
                    CloneURL string `json:"clone_url"`
                } `json:"repo"`
            } `json:"head"`
        } `json:"pull_request"`
    }

    json.Unmarshal(raw, &payload)

    if payload.Action != "submitted" {
        return
    }

    CreateJob(apiKey, "github_pr_apply", map[string]string{
        "repo":       payload.PullRequest.Head.Repo.CloneURL,
        "branch":     payload.PullRequest.Head.Ref,
        "project_id": fmtInt(projectID),
    })
}
