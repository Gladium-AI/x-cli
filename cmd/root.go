package cmd

import (
	"fmt"
	"os"

	"github.com/paolo/x-cli/internal/api"
	"github.com/spf13/cobra"
)

var jsonOutput bool
var verboseOutput bool

var rootCmd = &cobra.Command{
	Use:   "x-cli",
	Short: "A CLI tool for interacting with X (Twitter)",
	Long:  "x-cli wraps X's internal GraphQL API to fetch timelines, tweets, users, and more.",
}

func Execute() {
	// Wire verbose flag before any command runs
	cobra.OnInitialize(func() {
		api.VerboseMode = verboseOutput
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output raw JSON")
	rootCmd.PersistentFlags().BoolVar(&verboseOutput, "verbose", false, "Print request URLs and response details")
}
