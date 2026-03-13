package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/paolo/x-cli/internal/api"
	"github.com/paolo/x-cli/internal/output"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User profile commands",
}

var userGetCmd = &cobra.Command{
	Use:   "get <handle>",
	Short: "Get a user's profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handle := strings.TrimPrefix(args[0], "@")
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		user, rawJSON, err := client.GetUserByScreenName(context.Background(), handle)
		if err != nil {
			return fmt.Errorf("fetch user @%s: %w", handle, err)
		}

		output.PrintUser(user, jsonOutput, rawJSON)
		return nil
	},
}

func init() {
	userCmd.AddCommand(userGetCmd)
	rootCmd.AddCommand(userCmd)
}
