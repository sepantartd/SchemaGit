package store

import "github.com/google/uuid"

type Webhook struct {
	ID   string `json:"id"`
	OrgID string `json:"org_id"`
	URL  string `json:"url"`
}

func AddWebhook(orgID, rawURL string) error {
	_, err := DB.Exec(`
		INSERT INTO webhooks (id, org_id, url)
		VALUES (?, ?, ?)
	`, uuid.New().String(), orgID, rawURL)
	return err
}

func ListWebhooks(orgID string) ([]Webhook, error) {
	rows, err := DB.Query(`
		SELECT id, org_id, url FROM webhooks
		WHERE org_id = ? ORDER BY id
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	webhooks := []Webhook{}
	for rows.Next() {
		var webhook Webhook
		if err := rows.Scan(&webhook.ID, &webhook.OrgID, &webhook.URL); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func DeleteWebhook(orgID, id string) error {
	_, err := DB.Exec(`DELETE FROM webhooks WHERE org_id = ? AND id = ?`, orgID, id)
	return err
}
