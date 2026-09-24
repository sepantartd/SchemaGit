package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Execute migration via Cloud pipeline",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Project name required")
            return
        }

        project := args[0]
        token := os.Getenv("SCHEMAGIT_TOKEN")

        body := map[string]string{"project": project}
        buf, _ := json.Marshal(body)

        req, _ := http.NewRequest("POST", cloudURL+"/api/cli/migrate", bytes.NewBuffer(buf))
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        fmt.Println("Migration Job:")
        fmt.Println(data["job"])
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
