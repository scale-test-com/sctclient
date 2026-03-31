package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"scale-test/cli/internal/model"
)

var (
	createScenarioID   int
	createFile         string
	createWait         bool
	createPollInterval time.Duration
)

// terminalStates are run states that indicate the run has finished.
var terminalStates = map[string]bool{
	"completed": true,
	"failed":    true,
	"cancelled": true,
}

var runCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and start a new load test run",
	Long: `Create and start a new load test run.

Provide the scenario via one of these flags:
  --scenario-id <int>   Use a saved scenario by its ID
  --file <path>         Load a YAML scenario file

Both flags may be combined; the server gives priority to the inline content.
Use --wait to poll until the run reaches a terminal state (completed/failed/cancelled).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reqBody := model.CreateRunRequest{}

		if createFile != "" {
			data, err := os.ReadFile(createFile)
			if err != nil {
				return fmt.Errorf("read file %s: %w", createFile, err)
			}
			var content model.ScenarioContent
			if err := yaml.Unmarshal(data, &content); err != nil {
				return fmt.Errorf("parse YAML scenario: %w", err)
			}
			reqBody.Content = &content
		}

		if createScenarioID != 0 {
			id := createScenarioID
			reqBody.ScenarioID = &id
		}

		if reqBody.Content == nil && reqBody.ScenarioID == nil {
			return fmt.Errorf("provide at least one of --scenario-id or --file")
		}

		resp, err := client.CreateRun(reqBody)
		if err != nil {
			return err
		}

		if !createWait {
			return printJSON(resp)
		}

		fmt.Fprintf(os.Stderr, "Run %s created (state: %s). Waiting for completion...\n", resp.ID, resp.State)

		for {
			run, err := client.GetRun(resp.ID)
			if err != nil {
				return fmt.Errorf("poll run status: %w", err)
			}
			fmt.Fprintf(os.Stderr, "  state: %s\n", run.State)
			if terminalStates[run.State] {
				return printJSON(run)
			}
			time.Sleep(createPollInterval)
		}
	},
}

func init() {
	runCreateCmd.Flags().IntVar(&createScenarioID, "scenario-id", 0, "ID of a saved scenario")
	runCreateCmd.Flags().StringVar(&createFile, "file", "", "Path to a YAML scenario file")
	runCreateCmd.Flags().BoolVar(&createWait, "wait", false, "Wait until the run reaches a terminal state")
	runCreateCmd.Flags().DurationVar(&createPollInterval, "poll-interval", 2*time.Second, "Polling interval when --wait is set")
	runCmd.AddCommand(runCreateCmd)
}

// printJSON encodes v as indented JSON and writes it to stdout.
func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
