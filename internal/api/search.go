package api

import (
	"context"

	"github.com/paolo/x-cli/internal/models"
	"github.com/tidwall/gjson"
)

// SearchProducts maps CLI search type flags to X API product values.
var SearchProducts = map[string]string{
	"top":    "Top",
	"latest": "Latest",
	"people": "People",
	"media":  "Photos",
}

// Search performs a search query.
func (c *Client) Search(ctx context.Context, query string, product string, count int, cursor string) (models.TimelinePage, []byte, error) {
	vars := map[string]interface{}{
		"rawQuery":               query,
		"count":                  count,
		"querySource":            "typed_query",
		"product":                product,
		"withGrokTranslatedBio":  false,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	data, err := c.GraphQLGet(ctx, Endpoints["SearchTimeline"], vars)
	if err != nil {
		return models.TimelinePage{}, nil, err
	}

	instructions := gjson.GetBytes(data, "data.search_by_raw_query.search_timeline.timeline.instructions")
	page := models.ParseTimeline(instructions)
	return page, data, nil
}
