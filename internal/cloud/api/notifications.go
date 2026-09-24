package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/cloud/store"
)

func NotificationSetHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    orgs, err := store.UserOrgs(email)
    if err != nil {
        http.Error(w, "unable to load organizations", http.StatusInternalServerError)
        return
    }
    if len(orgs) == 0 {
        http.Error(w, "no org found", http.StatusBadRequest)
        return
    }

    var body struct {
        Email        bool   `json:"email"`
        Slack        bool   `json:"slack"`
        SlackWebhook string `json:"slack_webhook"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    orgID := orgs[0].ID
    if err := store.SetNotification(orgID, body.Email, body.Slack, body.SlackWebhook); err != nil {
        http.Error(w, "unable to save notification settings", http.StatusInternalServerError)
        return
    }
    store.AddAudit(orgID, "notification.update", "updated settings")

    _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func NotificationGetHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    orgs, err := store.UserOrgs(email)
    if err != nil {
        http.Error(w, "unable to load organizations", http.StatusInternalServerError)
        return
    }
    if len(orgs) == 0 {
        http.Error(w, "no org found", http.StatusBadRequest)
        return
    }

    s, err := store.GetNotification(orgs[0].ID)
    if err != nil && s == nil {
        s = &store.NotificationSetting{OrgID: orgs[0].ID}
    }
    _ = json.NewEncoder(w).Encode(s)
}
