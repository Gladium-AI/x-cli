package api

import (
	"context"

	"github.com/paolo/x-cli/internal/models"
	"github.com/tidwall/gjson"
)

// GetUserByScreenName fetches a user profile by handle.
func (c *Client) GetUserByScreenName(ctx context.Context, screenName string) (models.User, []byte, error) {
	vars := map[string]interface{}{
		"screen_name":            screenName,
		"withGrokTranslatedBio":  false,
	}

	data, err := c.GraphQL(ctx, Endpoints["UserByScreenName"], vars)
	if err != nil {
		return models.User{}, nil, err
	}

	result := gjson.GetBytes(data, "data.user.result")
	user := models.ParseUser(result)
	return user, data, nil
}

// ResolveUserID resolves a screen name to a user rest_id.
func (c *Client) ResolveUserID(ctx context.Context, screenName string) (string, error) {
	user, _, err := c.GetUserByScreenName(ctx, screenName)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}
