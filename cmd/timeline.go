package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var timelineCount int
var timelineCursor string

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

		// Resolve handle to user ID
		userID, err := client.ResolveUserID(context.Background(), handle)
		if err != nil {
			return fmt.Errorf("resolve @%s: %w", handle, err)
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

func init() {
	timelineCmd.PersistentFlags().IntVar(&timelineCount, "count", 20, "Number of tweets to fetch")
	timelineCmd.PersistentFlags().StringVar(&timelineCursor, "cursor", "", "Pagination cursor")

	timelineCmd.AddCommand(timelineHomeCmd, timelineUserCmd)
	rootCmd.AddCommand(timelineCmd)
}
