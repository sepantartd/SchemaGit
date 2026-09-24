package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <project>",
	Short: "Generate schema diff from Cloud",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := cliRequest(cmd, http.MethodGet, "/api/cli/diff", args[0], nil)
		if err != nil {
			return err
		}

		value, err := responseField(body, "diff")
		if err != nil {
			return err
		}
		fmt.Println("Diff:")
		fmt.Println(value)
		return nil
	},
}

func cliRequest(cmd *cobra.Command, method, endpoint, project string, payload any) ([]byte, error) {
	token := os.Getenv("SCHEMAGIT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("SCHEMAGIT_TOKEN is not set")
	}
	if strings.TrimSpace(cloudURL) == "" {
		return nil, fmt.Errorf("cloud URL is not configured")
	}

	base, err := url.Parse(cloudURL)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud URL: %w", err)
	}
	base.Path = strings.TrimRight(base.Path, "/") + endpoint
	query := base.Query()
	query.Set("project", project)
	base.RawQuery = query.Encode()

	var reader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(cmd.Context(), method, base.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-Token", token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(string(data))
		if message == "" {
			message = res.Status
		}
		return nil, fmt.Errorf("cloud API returned %s: %s", res.Status, message)
	}
	return data, nil
}

func responseField(body []byte, field string) (string, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("decode cloud response: %w", err)
	}
	raw, ok := data[field]
	if !ok || string(raw) == "null" {
		return "", fmt.Errorf("cloud response is missing %q", field)
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text, nil
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		return "", fmt.Errorf("invalid %q in cloud response: %w", field, err)
	}
	return pretty.String(), nil
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
