package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <project> <output-file>",
	Short: "Pull schema from Cloud",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := cliRequest(cmd, http.MethodGet, "/api/cli/pull", args[0], nil)
		if err != nil {
			return err
		}
		schema, err := responseField(body, "schema")
		if err != nil {
			return err
		}
		if err := os.WriteFile(args[1], []byte(schema), 0644); err != nil {
			return fmt.Errorf("write schema file: %w", err)
		}
		fmt.Println("Schema saved to", args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
