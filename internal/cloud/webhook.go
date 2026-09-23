package cloud

import (
    "encoding/json"
    "net/http"
)

func StartWebhookServer(apiKey string) {
    http.HandleFunc("/webhook/github", func(w http.ResponseWriter, r *http.Request) {

        var payload map[string]interface{}
        json.NewDecoder(r.Body).Decode(&payload)

        repo := payload["repository"].(map[string]interface{})["clone_url"].(string)
        branch := payload["ref"].(string)

        projectID := findProjectByWebhook(r.URL.Query().Get("id"))

        CreateJob(apiKey, "github_push", map[string]string{
            "repo":       repo,
            "branch":     branch,
            "project_id": fmtInt(projectID),
        })

        w.WriteHeader(200)
    })

    http.ListenAndServe(":8090", nil)
}
