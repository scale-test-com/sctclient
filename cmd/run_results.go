package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runResultsCmd = &cobra.Command{
	Use:   "results <run-uuid>",
	Short: "Get the aggregated results of a finished load test run",
	Long: `Get the aggregated results of a run: total requests, success rate, average
duration, requests per second (average and peak) and HTTP status code breakdown.

Results are available once the run is finished. While it is still in progress
the API answers 409.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := client.GetRunResults(args[0])
		if err != nil {
			return fmt.Errorf("get run results: %w", err)
		}
		return printJSON(results)
	},
}

func init() {
	runCmd.AddCommand(runResultsCmd)
}
