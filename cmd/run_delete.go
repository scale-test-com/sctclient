package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runDeleteCmd = &cobra.Command{
	Use:   "delete <run-uuid>",
	Short: "Delete a load test run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := client.DeleteRun(args[0])
		if err != nil {
			return fmt.Errorf("delete run: %w", err)
		}
		return printJSON(resp)
	},
}

func init() {
	runCmd.AddCommand(runDeleteCmd)
}
