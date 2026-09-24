package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan <project>",
	Short: "Generate migration plan from Cloud",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := cliRequest(cmd, http.MethodGet, "/api/cli/plan", args[0], nil)
		if err != nil {
			return err
		}
		value, err := responseField(body, "plan")
		if err != nil {
			return err
		}
		fmt.Println("Plan:")
		fmt.Println(value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
