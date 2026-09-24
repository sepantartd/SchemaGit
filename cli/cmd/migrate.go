package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate <project>",
	Short: "Execute migration via Cloud pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := cliRequest(cmd, http.MethodPost, "/api/cli/migrate", args[0], map[string]string{"project": args[0]})
		if err != nil {
			return err
		}
		value, err := responseField(body, "job")
		if err != nil {
			return err
		}
		fmt.Println("Migration Job:")
		fmt.Println(value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
