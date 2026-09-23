package api

import (
    "encoding/json"
    "net/http"
    "fmt"
)

type GitHubWebhook struct {
    Action     string `json:"action"`
    PullRequest struct {
        Number int    `json:"number"`
        Title  string `json:"title"`
        Head   struct {
            Sha string `json:"sha"`
        } `json:"head"`
    } `json:"pull_request"`
    Repository struct {
        FullName string `json:"full_name"`
    } `json:"repository"`
}

func GitHubWebhookHandler(w http.ResponseWriter, req *http.Request) {
    var payload GitHubWebhook
    json.NewDecoder(req.Body).Decode(&payload)

    fmt.Println("Webhook received:", payload.Action)

    switch payload.Action {

    case "opened":
        fmt.Printf("PR #%d opened: %s\n", payload.PullRequest.Number, payload.PullRequest.Title)
        // TODO: Trigger diff + plan + store result

    case "synchronize":
        fmt.Printf("PR #%d updated (new commit: %s)\n",
            payload.PullRequest.Number,
            payload.PullRequest.Head.Sha,
        )
        // TODO: Trigger diff + plan + update result

    case "closed":
        fmt.Printf("PR #%d closed\n", payload.PullRequest.Number)
        // TODO: Cleanup or finalize

    default:
        fmt.Println("Unhandled webhook action:", payload.Action)
    }

    w.WriteHeader(http.StatusOK)
}
