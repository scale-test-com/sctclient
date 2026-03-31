package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"scale-test/cli/internal/api"
)

var (
	apiKey  string
	baseURL string
	client  *api.Client
)

var rootCmd = &cobra.Command{
	Use:   "scale-test",
	Short: "CLI client for the Scale-Test load testing service",
	Long: `scale-test is a command-line client for the Scale-Test load testing API.

Authentication:
  --api-key flag  or  SCALE_TEST_API_KEY environment variable

Base URL override:
  --base-url flag  or  SCALE_TEST_BASE_URL environment variable
  (default: https://scale-test.com/api/v1)`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// api-key: flag > env
		if apiKey == "" {
			apiKey = os.Getenv("SCALE_TEST_API_KEY")
		}
		if apiKey == "" {
			return fmt.Errorf("API key required: use --api-key flag or set SCALE_TEST_API_KEY")
		}

		// base-url: flag > env > hardcoded default
		if baseURL == "" {
			if v := os.Getenv("SCALE_TEST_BASE_URL"); v != "" {
				baseURL = v
			} else {
				baseURL = "https://scale-test.com/api/v1"
			}
		}

		client = api.NewClient(baseURL, apiKey)
		return nil
	},
}

// Execute runs the root cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (env: SCALE_TEST_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "API base URL (env: SCALE_TEST_BASE_URL)")
}
