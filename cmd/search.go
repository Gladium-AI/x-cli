package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var searchCount int
var searchCursor string
var searchType string

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for tweets",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		product, ok := api.SearchProducts[searchType]
		if !ok {
			return fmt.Errorf("invalid search type: %s (use: top, latest, people, media)", searchType)
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		page, rawJSON, err := client.Search(context.Background(), query, product, searchCount, searchCursor)
		if err != nil {
			return fmt.Errorf("search '%s': %w", query, err)
		}

		output.PrintTweets(page.Tweets, jsonOutput, rawJSON)
		if !jsonOutput {
			output.PrintCursor(fmt.Sprintf("x-cli search \"%s\"", query), page.NextCursor)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchCount, "count", 20, "Number of results")
	searchCmd.Flags().StringVar(&searchCursor, "cursor", "", "Pagination cursor")
	searchCmd.Flags().StringVar(&searchType, "type", "top", "Search type: top, latest, people, media")
	rootCmd.AddCommand(searchCmd)
}
