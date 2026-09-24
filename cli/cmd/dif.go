Copilot, please create the following new CLI command files for SchemaGit.
Each file should contain exactly the code provided below.

-----------------------------------------
FILE: cli/cmd/diff.go
-----------------------------------------
package cmd

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
    Use:   "diff",
    Short: "Generate schema diff from Cloud",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Project name required")
            return
        }

        project := args[0]
        token := os.Getenv("SCHEMAGIT_TOKEN")

        req, _ := http.NewRequest("GET", cloudURL+"/api/cli/diff?project="+project, nil)
        req.Header.Set("X-API-Token", token)

        res, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var data map[string]interface{}
        json.NewDecoder(res.Body).Decode(&data)

        fmt.Println("Diff:")
        fmt.Println(data["diff"])
    },
}

func init() {
    rootCmd.AddCommand(diffCmd)
}

-----------------------------------------
FILE: cli/cmd/plan.go
-----------------------------------------
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

-----------------------------------------
FILE: cli/cmd/migrate.go
-----------------------------------------
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

-----------------------------------------
FILE: cli/cmd/push.go
-----------------------------------------
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

        _, err = http.DefaultClient.Do(req)
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

-----------------------------------------
FILE: cli/cmd/pull.go
-----------------------------------------
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
