package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paolo/x-cli/internal/models"
)

// PrintTweets displays a list of tweets in the selected format.
func PrintTweets(tweets []*models.Tweet, jsonMode bool, rawJSON []byte) {
	if jsonMode {
		printRawJSON(rawJSON)
		return
	}
	for i, t := range tweets {
		PrintTweet(t)
		if i < len(tweets)-1 {
			fmt.Println()
		}
	}
}

// PrintTweet displays a single tweet.
func PrintTweet(t *models.Tweet) {
	printTweetPretty(t)
}

// PrintUser displays a user profile in the selected format.
func PrintUser(u models.User, jsonMode bool, rawJSON []byte) {
	if jsonMode {
		printRawJSON(rawJSON)
		return
	}
	printUserPretty(u)
}

// PrintUsers displays a list of users.
func PrintUsers(users []models.User, jsonMode bool, rawJSON []byte) {
	if jsonMode {
		printRawJSON(rawJSON)
		return
	}
	for i, u := range users {
		printUserPretty(u)
		if i < len(users)-1 {
			fmt.Println()
		}
	}
}

// PrintCursor shows the pagination hint.
func PrintCursor(command, cursor string) {
	if cursor == "" {
		return
	}
	fmt.Printf("\nNext page: %s --cursor \"%s\"\n", command, cursor)
}

func printRawJSON(data []byte) {
	// Pretty-print the JSON
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		os.Stdout.Write(data)
		return
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
