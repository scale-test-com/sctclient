package cmd

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage load test runs",
}

func init() {
	rootCmd.AddCommand(runCmd)
}
