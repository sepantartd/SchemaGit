package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
    Use:   "pull",
    Short: "Pull schema from Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 2 {
            fmt.Println("Usage: schemagit pull <project> <output-file>")
            return
        }

        project := args[0]
        out := args[1]

        token := os.Getenv("SCHEMAGIT_TOKEN")

        req, _ := http.NewRequest("GET", cloudURL+"/api/cli/pull?project="+project, nil)
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        os.WriteFile(out, []byte(data["schema"].(string)), 0644)

        fmt.Println("Schema saved to", out)
    },
}

func init() {
    rootCmd.AddCommand(pullCmd)
}
