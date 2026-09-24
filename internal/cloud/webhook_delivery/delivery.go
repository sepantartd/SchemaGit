package webhook_delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sepanta/schemagit/internal/cloud/store"
)

const maxResponseBody = 1 << 20

var client = &http.Client{Timeout: 10 * time.Second}

type envelope struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

func ValidateURL(raw string) error {
	u, err := url.ParseRequestURI(strings.TrimSpace(raw))
	if err != nil || u == nil || u.Host == "" || u.User != nil {
		return fmt.Errorf("webhook URL must be an absolute HTTP or HTTPS URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https")
	}
	return nil
}

func Dispatch(orgID, event string, payload interface{}) {
	webhooks, err := store.ListWebhooks(orgID)
	if err != nil {
		return
	}
	body, err := json.Marshal(envelope{Event: event, Payload: payload})
	if err != nil {
		return
	}
	for _, webhook := range webhooks {
		if err := deliver(webhook.URL, body); err != nil {
			// Delivery is best effort: one endpoint must not block the others.
			continue
		}
	}
}

func deliver(rawURL string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = res.Body.Read(make([]byte, maxResponseBody))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("webhook returned %s", res.Status)
	}
	return nil
}
