package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/models"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var timelineCount int
var timelineCursor string
var timelineAll bool
var timelineMaxPages int

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Timeline commands",
}

var timelineHomeCmd = &cobra.Command{
	Use:   "home",
	Short: "View your home timeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		if timelineAll {
			return paginateTimeline(client, "x-cli timeline home", func(cursor string) (models.TimelinePage, []byte, error) {
				return client.GetHomeTimeline(context.Background(), timelineCount, cursor)
			})
		}

		page, rawJSON, err := client.GetHomeTimeline(context.Background(), timelineCount, timelineCursor)
		if err != nil {
			return fmt.Errorf("fetch home timeline: %w", err)
		}

		output.PrintTweets(page.Tweets, jsonOutput, rawJSON)
		if !jsonOutput {
			output.PrintCursor("x-cli timeline home", page.NextCursor)
		}
		return nil
	},
}

var timelineUserCmd = &cobra.Command{
	Use:   "user <handle>",
	Short: "View a user's tweets",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handle := strings.TrimPrefix(args[0], "@")
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		userID, err := client.ResolveUserID(context.Background(), handle)
		if err != nil {
			return fmt.Errorf("resolve @%s: %w", handle, err)
		}

		if timelineAll {
			return paginateTimeline(client, fmt.Sprintf("x-cli timeline user @%s", handle), func(cursor string) (models.TimelinePage, []byte, error) {
				return client.GetUserTweets(context.Background(), userID, timelineCount, cursor)
			})
		}

		page, rawJSON, err := client.GetUserTweets(context.Background(), userID, timelineCount, timelineCursor)
		if err != nil {
			return fmt.Errorf("fetch tweets for @%s: %w", handle, err)
		}

		output.PrintTweets(page.Tweets, jsonOutput, rawJSON)
		if !jsonOutput {
			output.PrintCursor(fmt.Sprintf("x-cli timeline user @%s", handle), page.NextCursor)
		}
		return nil
	},
}

type timelineFetcher func(cursor string) (models.TimelinePage, []byte, error)

func paginateTimeline(client *api.Client, cmdName string, fetch timelineFetcher) error {
	cursor := timelineCursor
	totalTweets := 0

	for page := 0; page < timelineMaxPages; page++ {
		result, rawJSON, err := fetch(cursor)
		if err != nil {
			return err
		}

		if jsonOutput {
			output.PrintTweets(nil, true, rawJSON)
		} else {
			output.PrintTweets(result.Tweets, false, nil)
		}

		totalTweets += len(result.Tweets)
		cursor = result.NextCursor

		if cursor == "" || len(result.Tweets) == 0 {
			break
		}

		// Check rate limits between pages
		if client.LastRateLimit != nil {
			client.LastRateLimit.WaitIfNeeded()
		}
	}

	if !jsonOutput {
		fmt.Printf("\nFetched %d tweets across %d pages.\n", totalTweets, min(timelineMaxPages, totalTweets/timelineCount+1))
	}
	return nil
}

func init() {
	timelineCmd.PersistentFlags().IntVar(&timelineCount, "count", 20, "Number of tweets per page")
	timelineCmd.PersistentFlags().StringVar(&timelineCursor, "cursor", "", "Pagination cursor")
	timelineCmd.PersistentFlags().BoolVar(&timelineAll, "all", false, "Auto-paginate through all results")
	timelineCmd.PersistentFlags().IntVar(&timelineMaxPages, "max-pages", 10, "Max pages when using --all")

	timelineCmd.AddCommand(timelineHomeCmd, timelineUserCmd)
	rootCmd.AddCommand(timelineCmd)
}
