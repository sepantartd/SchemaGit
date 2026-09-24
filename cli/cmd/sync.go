package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
    Use:   "sync",
    Short: "Sync CLI with SchemaGit Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        token := os.Getenv("SCHEMAGIT_TOKEN")
        if token == "" {
            fmt.Println("Missing SCHEMAGIT_TOKEN")
            return
        }

        req, _ := http.NewRequest("GET", cloudURL+"/api/cli/sync", nil)
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        fmt.Println("Email:", data["email"])
        fmt.Println("Projects:", data["projects"])
    },
}

func init() {
    rootCmd.AddCommand(syncCmd)
}
