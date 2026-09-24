package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cloudURL = strings.TrimRight(os.Getenv("SCHEMAGIT_CLOUD_URL"), "/")

var rootCmd = &cobra.Command{
	Use:           "schemagit",
	Short:         "SchemaGit CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	if cloudURL == "" {
		cloudURL = "http://localhost:9090"
	}
	rootCmd.PersistentFlags().StringVar(&cloudURL, "cloud-url", cloudURL, "SchemaGit Cloud URL")
}

// Execute runs the SchemaGit Cobra command tree.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
