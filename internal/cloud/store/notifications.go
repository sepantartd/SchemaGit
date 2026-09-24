package store

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/smtp"
    "os"
)

type NotificationSetting struct {
    OrgID       string `json:"org_id"`
    Email       bool   `json:"email"`
    Slack       bool   `json:"slack"`
    SlackWebhook string `json:"slack_webhook"`
}

func SetNotification(orgID string, email bool, slack bool, slackWebhook string) error {
    _, err := DB.Exec(`
        INSERT OR REPLACE INTO notifications (org_id, email, slack, slack_webhook)
        VALUES (?, ?, ?, ?)
    `, orgID, email, slack, slackWebhook)
    return err
}

func GetNotification(orgID string) (*NotificationSetting, error) {
    row := DB.QueryRow(`
        SELECT org_id, email, slack, slack_webhook
        FROM notifications WHERE org_id = ?
    `, orgID)

    var s NotificationSetting
    err := row.Scan(&s.OrgID, &s.Email, &s.Slack, &s.SlackWebhook)
    if err != nil {
        return nil, err
    }

    return &s, nil
}

func OrgOwner(orgID string) (string, error) {
    row := DB.QueryRow(`SELECT owner FROM orgs WHERE id = ?`, orgID)
    var owner string
    if err := row.Scan(&owner); err != nil {
        return "", err
    }
    return owner, nil
}

func SendEmail(to string, subject string, body string) {
    if to == "" {
        return
    }
    host := os.Getenv("SMTP_HOST")
    user := os.Getenv("SMTP_USER")
    pass := os.Getenv("SMTP_PASS")
    if host == "" || user == "" {
        return
    }

    msg := fmt.Sprintf("Subject: %s\n\n%s", subject, body)
    _ = smtp.SendMail(host, smtp.PlainAuth("", user, pass, host), user, []string{to}, []byte(msg))
}

func SendSlack(webhook string, text string) {
    payload := map[string]string{"text": text}
    buf, err := json.Marshal(payload)
    if err != nil {
        return
    }
    req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err == nil && resp.Body != nil {
        resp.Body.Close()
    }
}

func DispatchNotification(orgID string, event string, payload string, email string) {
    s, err := GetNotification(orgID)
    if err != nil {
        return
    }

    if s.Email {
        SendEmail(email, "SchemaGit Notification: "+event, payload)
    }

    if s.Slack && s.SlackWebhook != "" {
        SendSlack(s.SlackWebhook, fmt.Sprintf("%s: %s", event, payload))
    }
}
