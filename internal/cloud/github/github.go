package github

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "schemagit/internal/schema/diff"
)

type GitHubClient struct {
    Token string
}

func NewClient(token string) *GitHubClient {
    return &GitHubClient{Token: token}
}

// گرفتن فایل از GitHub
func (c *GitHubClient) FetchFile(repo string, commitSHA string, path string) ([]byte, error) {
    url := fmt.Sprintf(
        "https://raw.githubusercontent.com/%s/%s/%s",
        repo,
        commitSHA,
        path,
    )

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+c.Token)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    return ioutil.ReadAll(res.Body)
}

// گرفتن اسکیماها از ریپو
func (c *GitHubClient) FetchSchemas(repo string, prNumber int, commitSHA string) ([]diff.Table, []diff.Table) {
    // مسیرهای اسکیما
    oldPath := "schema.old.json"
    newPath := "schema.new.json"

    oldData, _ := c.FetchFile(repo, commitSHA, oldPath)
    newData, _ := c.FetchFile(repo, commitSHA, newPath)

    var oldSchema []diff.Table
    var newSchema []diff.Table

    json.Unmarshal(oldData, &oldSchema)
    json.Unmarshal(newData, &newSchema)

    return oldSchema, newSchema
}
