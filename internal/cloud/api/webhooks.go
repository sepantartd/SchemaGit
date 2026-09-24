package api

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/sepanta/schemagit/internal/cloud/store"
    "github.com/sepanta/schemagit/internal/cloud/webhook_delivery"
)

func WebhookAddHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    email, err := authenticatedEmail(req)
    if err != nil {
        writeJSONError(w, http.StatusUnauthorized, "invalid token")
        return
    }
    var body struct {
        OrgID string `json:"org_id"`
        URL   string `json:"url"`
    }
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&body); err != nil {
        writeJSONError(w, http.StatusBadRequest, "invalid JSON request")
        return
    }
    if body.OrgID == "" || strings.TrimSpace(body.URL) == "" {
        writeJSONError(w, http.StatusBadRequest, "org_id and url are required")
        return
    }
    if !store.CanManageOrg(email, body.OrgID) {
        writeJSONError(w, http.StatusForbidden, "forbidden")
        return
    }
    body.URL = strings.TrimSpace(body.URL)
    if err := webhook_delivery.ValidateURL(body.URL); err != nil {
        writeJSONError(w, http.StatusBadRequest, err.Error())
        return
    }
    if err := store.AddWebhook(body.OrgID, body.URL); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "unable to add webhook")
        return
    }
    store.AddAudit(body.OrgID, "webhook.add", body.URL)
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func WebhookListHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
        writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    email, err := authenticatedEmail(req)
    if err != nil {
        writeJSONError(w, http.StatusUnauthorized, "invalid token")
        return
    }
    orgID := req.URL.Query().Get("org")
    if orgID == "" || !store.CanAccessOrg(email, orgID) {
        writeJSONError(w, http.StatusForbidden, "forbidden")
        return
    }
    webhooks, err := store.ListWebhooks(orgID)
    if err != nil {
        writeJSONError(w, http.StatusInternalServerError, "unable to list webhooks")
        return
    }
    writeJSON(w, http.StatusOK, webhooks)
}

func authenticatedEmail(req *http.Request) (string, error) {
    return store.ValidateToken(req.Header.Get("Authorization"))
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(value)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, map[string]string{"error": message})
}
