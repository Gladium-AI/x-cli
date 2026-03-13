package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var tweetURLPattern = regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:twitter\.com|x\.com)/\w+/status/(\d+)`)

var tweetCmd = &cobra.Command{
	Use:   "tweet",
	Short: "Tweet commands",
}

var tweetGetCmd = &cobra.Command{
	Use:   "get <tweet_id_or_url>",
	Short: "Get a tweet by ID or URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tweetID := parseTweetID(args[0])
		if tweetID == "" {
			return fmt.Errorf("invalid tweet ID or URL: %s", args[0])
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tweet, rawJSON, err := client.GetTweetDetail(context.Background(), tweetID)
		if err != nil {
			return fmt.Errorf("fetch tweet %s: %w", tweetID, err)
		}

		if jsonOutput {
			output.PrintTweets(nil, true, rawJSON)
		} else {
			output.PrintTweet(tweet)
		}
		return nil
	},
}

func parseTweetID(input string) string {
	// Try URL match first
	if matches := tweetURLPattern.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1]
	}
	// Check if it's a plain numeric ID
	input = strings.TrimSpace(input)
	for _, c := range input {
		if c < '0' || c > '9' {
			return ""
		}
	}
	if len(input) > 0 {
		return input
	}
	return ""
}

func init() {
	tweetCmd.AddCommand(tweetGetCmd)
	rootCmd.AddCommand(tweetCmd)
}
