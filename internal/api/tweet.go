package api

import (
	"context"
	"fmt"

	"github.com/paolo/x-cli/internal/models"
	"github.com/tidwall/gjson"
)

// GetTweetDetail fetches a single tweet by ID.
func (c *Client) GetTweetDetail(ctx context.Context, tweetID string) (*models.Tweet, []byte, error) {
	vars := map[string]interface{}{
		"focalTweetId":                           tweetID,
		"with_rux_injections":                    false,
		"rankingMode":                            "Relevance",
		"includePromotedContent":                 true,
		"withCommunity":                          true,
		"withQuickPromoteEligibilityTweetFields": true,
		"withBirdwatchNotes":                     true,
		"withVoice":                              true,
	}

	data, err := c.GraphQL(ctx, Endpoints["TweetDetail"], vars)
	if err != nil {
		return nil, nil, err
	}

	// The focal tweet is in the instructions entries with entryId matching "tweet-<id>"
	instructions := gjson.GetBytes(data, "data.threaded_conversation_with_injections_v2.instructions")
	targetEntryID := "tweet-" + tweetID

	var tweet *models.Tweet
	instructions.ForEach(func(_, instr gjson.Result) bool {
		instr.Get("entries").ForEach(func(_, entry gjson.Result) bool {
			if entry.Get("entryId").String() == targetEntryID {
				result := entry.Get("content.itemContent.tweet_results.result")
				tweet = models.ParseTweet(result)
				return false // found it
			}
			return true
		})
		return tweet == nil
	})

	if tweet == nil {
		return nil, data, fmt.Errorf("tweet %s not found in response", tweetID)
	}

	return tweet, data, nil
}
