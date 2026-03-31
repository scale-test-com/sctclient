package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runGetCmd = &cobra.Command{
	Use:   "get <run-uuid>",
	Short: "Get details of a load test run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run, err := client.GetRun(args[0])
		if err != nil {
			return fmt.Errorf("get run: %w", err)
		}
		return printJSON(run)
	},
}

func init() {
	runCmd.AddCommand(runGetCmd)
}
