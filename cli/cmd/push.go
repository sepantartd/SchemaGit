package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
    Use:   "push",
    Short: "Push local schema to Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 2 {
            fmt.Println("Usage: schemagit push <project> <schema-file>")
            return
        }

        project := args[0]
        file := args[1]

        schema, err := os.ReadFile(file)
        if err != nil {
            fmt.Println("Error reading schema:", err)
            return
        }

        token := os.Getenv("SCHEMAGIT_TOKEN")

        body := map[string]string{
            "project": project,
            "schema":  string(schema),
        }
        buf, _ := json.Marshal(body)

        req, _ := http.NewRequest("POST", cloudURL+"/api/cli/push", bytes.NewBuffer(buf))
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        fmt.Println("Schema pushed successfully")
    },
}

func init() {
    rootCmd.AddCommand(pushCmd)
}
