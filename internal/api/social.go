package api

import (
	"context"

	"github.com/paolo/x-cli/internal/models"
	"github.com/tidwall/gjson"
)

// GetFollowers fetches the followers list for a user.
func (c *Client) GetFollowers(ctx context.Context, userID string, count int, cursor string) ([]models.User, string, []byte, error) {
	return c.getSocialList(ctx, "Followers", userID, count, cursor)
}

// GetFollowing fetches the following list for a user.
func (c *Client) GetFollowing(ctx context.Context, userID string, count int, cursor string) ([]models.User, string, []byte, error) {
	return c.getSocialList(ctx, "Following", userID, count, cursor)
}

func (c *Client) getSocialList(ctx context.Context, endpoint string, userID string, count int, cursor string) ([]models.User, string, []byte, error) {
	vars := map[string]interface{}{
		"userId":                 userID,
		"count":                  count,
		"includePromotedContent": false,
		"withGrokTranslatedBio":  false,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	data, err := c.GraphQLGet(ctx, Endpoints[endpoint], vars)
	if err != nil {
		return nil, "", nil, err
	}

	// Social list responses have the same timeline structure
	instructions := gjson.GetBytes(data, "data.user.result.timeline.timeline.instructions")

	var users []models.User
	var nextCursor string

	instructions.ForEach(func(_, instr gjson.Result) bool {
		instr.Get("entries").ForEach(func(_, entry gjson.Result) bool {
			content := entry.Get("content")
			entryType := content.Get("entryType").String()

			switch entryType {
			case "TimelineTimelineItem":
				userResult := content.Get("itemContent.user_results.result")
				if userResult.Exists() {
					users = append(users, models.ParseUser(userResult))
				}
			case "TimelineTimelineCursor":
				if content.Get("cursorType").String() == "Bottom" {
					nextCursor = content.Get("value").String()
				}
			}
			return true
		})
		return true
	})

	return users, nextCursor, data, nil
}
