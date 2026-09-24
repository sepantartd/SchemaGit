package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var syncProjectCmd = &cobra.Command{
    Use:   "sync-project",
    Short: "Sync a specific project from Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Project name required")
            return
        }

        project := args[0]
        token := os.Getenv("SCHEMAGIT_TOKEN")

        req, _ := http.NewRequest("GET", cloudURL+"/api/cli/sync/project?project="+project, nil)
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        fmt.Println("Schema:")
        fmt.Println(data["schema"])

        fmt.Println("\nMigrations:")
        fmt.Println(data["migrations"])
    },
}

func init() {
    rootCmd.AddCommand(syncProjectCmd)
}
