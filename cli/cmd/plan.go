package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
    Use:   "plan",
    Short: "Generate migration plan from Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Project name required")
            return
        }

        project := args[0]
        token := os.Getenv("SCHEMAGIT_TOKEN")

        req, _ := http.NewRequest("GET", cloudURL+"/api/cli/plan?project="+project, nil)
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        fmt.Println("Plan:")
        fmt.Println(data["plan"])
    },
}

func init() {
    rootCmd.AddCommand(planCmd)
}
