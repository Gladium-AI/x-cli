package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/paolo/x-cli/internal/models"
)

var (
	handleColor = color.New(color.FgCyan, color.Bold)
	nameColor   = color.New(color.FgWhite, color.Bold)
	dimColor    = color.New(color.FgHiBlack)
	statColor   = color.New(color.FgHiBlack)
	separator   = dimColor.Sprint(strings.Repeat("─", 50))
)

func printTweetPretty(t *models.Tweet) {
	if t == nil {
		return
	}

	// Header: @handle · time ago
	if t.IsRetweet {
		dimColor.Print("♻ RT ")
	}
	handleColor.Printf("@%s", t.Author.ScreenName)
	if t.Author.IsBlueVerified {
		color.New(color.FgBlue).Print(" ✓")
	}
	dimColor.Printf(" · %s\n", t.TimeAgo())

	// Tweet text
	fmt.Println(t.FullText)

	// Media indicators
	for _, m := range t.Media {
		switch m.Type {
		case "photo":
			dimColor.Printf("  [image %dx%d] %s\n", m.Width, m.Height, m.URL)
		case "video":
			dimColor.Printf("  [video %dx%d] %s\n", m.Width, m.Height, m.VideoURL)
		case "animated_gif":
			dimColor.Printf("  [gif] %s\n", m.VideoURL)
		}
	}

	// Quoted tweet (compact)
	if t.QuotedTweet != nil {
		qt := t.QuotedTweet
		dimColor.Print("  ┌ ")
		handleColor.Printf("@%s", qt.Author.ScreenName)
		dimColor.Printf(": %s\n", truncateText(qt.FullText, 100))
	}

	// Stats line
	stats := []string{}
	if t.RetweetCount > 0 {
		stats = append(stats, fmt.Sprintf("♻ %s", models.FormatCount(t.RetweetCount)))
	}
	if t.LikeCount > 0 {
		stats = append(stats, fmt.Sprintf("♥ %s", models.FormatCount(t.LikeCount)))
	}
	if t.ReplyCount > 0 {
		stats = append(stats, fmt.Sprintf("💬 %s", models.FormatCount(t.ReplyCount)))
	}
	if t.ViewCount > 0 {
		stats = append(stats, fmt.Sprintf("👁 %s", models.FormatCount(t.ViewCount)))
	}
	if len(stats) > 0 {
		statColor.Println(strings.Join(stats, "  "))
	}

	fmt.Println(separator)
}

func printUserPretty(u models.User) {
	handleColor.Printf("@%s", u.ScreenName)
	if u.IsBlueVerified {
		color.New(color.FgBlue).Print(" ✓")
	}
	dimColor.Printf(" · ")
	nameColor.Println(u.Name)

	if u.Description != "" {
		fmt.Println(u.Description)
	}

	// Location + URL line
	var meta []string
	if u.Location != "" {
		meta = append(meta, "📍 "+u.Location)
	}
	if u.URL != "" {
		meta = append(meta, "🔗 "+u.URL)
	}
	if len(meta) > 0 {
		dimColor.Println(strings.Join(meta, " · "))
	}

	fmt.Printf("Followers: %s · Following: %s · Tweets: %s\n",
		models.FormatCount(u.FollowersCount),
		models.FormatCount(u.FriendsCount),
		models.FormatCount(u.StatusesCount),
	)

	if u.CreatedAt != "" {
		dimColor.Printf("Joined: %s\n", u.CreatedAt)
	}

	fmt.Println(separator)
}

func truncateText(s string, n int) string {
	// Replace newlines with spaces for compact display
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
