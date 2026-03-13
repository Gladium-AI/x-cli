package api

import (
	"context"

	"github.com/paolo/x-cli/internal/models"
	"github.com/tidwall/gjson"
)

// GetHomeTimeline fetches the authenticated user's home timeline.
func (c *Client) GetHomeTimeline(ctx context.Context, count int, cursor string) (models.TimelinePage, []byte, error) {
	vars := map[string]interface{}{
		"count":                  count,
		"includePromotedContent": true,
		"withCommunity":          true,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	data, err := c.GraphQLGet(ctx, Endpoints["HomeTimeline"], vars)
	if err != nil {
		return models.TimelinePage{}, nil, err
	}

	instructions := gjson.GetBytes(data, "data.home.home_timeline_urt.instructions")
	page := models.ParseTimeline(instructions)
	return page, data, nil
}

// GetUserTweets fetches tweets from a specific user.
func (c *Client) GetUserTweets(ctx context.Context, userID string, count int, cursor string) (models.TimelinePage, []byte, error) {
	vars := map[string]interface{}{
		"userId":                 userID,
		"count":                  count,
		"includePromotedContent": true,
		"withQuickPromoteEligibilityTweetFields": true,
		"withVoice": true,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	data, err := c.GraphQLGet(ctx, Endpoints["UserTweets"], vars)
	if err != nil {
		return models.TimelinePage{}, nil, err
	}

	instructions := gjson.GetBytes(data, "data.user.result.timeline_v2.timeline.instructions")
	page := models.ParseTimeline(instructions)
	return page, data, nil
}
