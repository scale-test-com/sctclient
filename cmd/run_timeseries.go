package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var timeseriesMetric string

var runTimeseriesCmd = &cobra.Command{
	Use:   "timeseries <run-uuid>",
	Short: "Get the per-second metrics of a finished load test run",
	Long: `Get the per-second metrics of a run (timestamps are UTC).

Use --metric to restrict the output to one series:
  status_codes, success_rate, average_duration, requests_per_second`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		series, err := client.GetRunTimeseries(args[0], timeseriesMetric)
		if err != nil {
			return fmt.Errorf("get run timeseries: %w", err)
		}
		return printJSON(series)
	},
}

func init() {
	runTimeseriesCmd.Flags().StringVar(&timeseriesMetric, "metric", "", "Only return this series (status_codes, success_rate, average_duration, requests_per_second)")
	runCmd.AddCommand(runTimeseriesCmd)
}
