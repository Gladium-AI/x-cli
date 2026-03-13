package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var socialCount int
var socialCursor string

var followersCmd = &cobra.Command{
	Use:   "followers <handle>",
	Short: "List a user's followers",
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

		users, nextCursor, rawJSON, err := client.GetFollowers(context.Background(), userID, socialCount, socialCursor)
		if err != nil {
			return fmt.Errorf("fetch followers for @%s: %w", handle, err)
		}

		output.PrintUsers(users, jsonOutput, rawJSON)
		if !jsonOutput {
			output.PrintCursor(fmt.Sprintf("x-cli followers @%s", handle), nextCursor)
		}
		return nil
	},
}

var followingCmd = &cobra.Command{
	Use:   "following <handle>",
	Short: "List who a user follows",
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

		users, nextCursor, rawJSON, err := client.GetFollowing(context.Background(), userID, socialCount, socialCursor)
		if err != nil {
			return fmt.Errorf("fetch following for @%s: %w", handle, err)
		}

		output.PrintUsers(users, jsonOutput, rawJSON)
		if !jsonOutput {
			output.PrintCursor(fmt.Sprintf("x-cli following @%s", handle), nextCursor)
		}
		return nil
	},
}

func init() {
	followersCmd.Flags().IntVar(&socialCount, "count", 20, "Number of users to fetch")
	followersCmd.Flags().StringVar(&socialCursor, "cursor", "", "Pagination cursor")

	followingCmd.Flags().IntVar(&socialCount, "count", 20, "Number of users to fetch")
	followingCmd.Flags().StringVar(&socialCursor, "cursor", "", "Pagination cursor")

	rootCmd.AddCommand(followersCmd, followingCmd)
}
