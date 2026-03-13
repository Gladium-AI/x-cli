package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/paolo/x-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to X via browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.BrowserLogin(context.Background())
		if err != nil {
			return err
		}
		if err := auth.Save(creds); err != nil {
			return err
		}
		color.Green("✓ Logged in as %s (user ID: %s)", creds.ScreenName, creds.UserID)
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current auth status",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Not logged in.")
			fmt.Fprintln(os.Stderr, "Run: x-cli auth login")
			return nil
		}
		color.Green("✓ Logged in")
		fmt.Printf("  Account:  %s\n", creds.ScreenName)
		fmt.Printf("  User ID:  %s\n", creds.UserID)
		fmt.Printf("  Since:    %s\n", creds.CreatedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.Clear(); err != nil {
			return err
		}
		fmt.Println("Logged out.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd, authStatusCmd, authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}
