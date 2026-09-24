package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <project> <schema-file>",
	Short: "Push local schema to Cloud",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		schema, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("read schema file: %w", err)
		}
		_, err = cliRequest(cmd, http.MethodPost, "/api/cli/push", args[0], map[string]string{
			"project": args[0],
			"schema":  string(schema),
		})
		if err != nil {
			return err
		}
		fmt.Println("Schema pushed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
