package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/paolo/x-cli/internal/auth"
)

const graphqlBaseURL = "https://x.com/i/api/graphql"

// Client handles authenticated requests to X's GraphQL API.
type Client struct {
	httpClient  *http.Client
	credentials *auth.Credentials
	LastRateLimit *RateLimit
}

// NewClient creates a new API client from stored credentials.
func NewClient() (*Client, error) {
	creds, err := auth.Load()
	if err != nil {
		return nil, err
	}
	return &Client{
		httpClient:  &http.Client{},
		credentials: creds,
	}, nil
}

// GraphQLGet performs a GET request to a GraphQL endpoint with variables encoded as query params.
func (c *Client) GraphQLGet(ctx context.Context, endpoint EndpointDef, variables interface{}) ([]byte, error) {
	varsJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("marshal variables: %w", err)
	}

	u := fmt.Sprintf("%s/%s/%s", graphqlBaseURL, endpoint.QueryID, endpoint.OperationName)
	params := url.Values{}
	params.Set("variables", string(varsJSON))

	if endpoint.HasFeatures {
		featJSON, _ := json.Marshal(DefaultFeatures)
		params.Set("features", string(featJSON))
	}
	if endpoint.HasFieldToggles {
		toggleJSON, _ := json.Marshal(DefaultFieldToggles)
		params.Set("fieldToggles", string(toggleJSON))
	}

	fullURL := u + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)
	return c.doRequest(req)
}

// GraphQLPost performs a POST request to a GraphQL endpoint with a JSON body.
func (c *Client) GraphQLPost(ctx context.Context, endpoint EndpointDef, variables interface{}) ([]byte, error) {
	u := fmt.Sprintf("%s/%s/%s", graphqlBaseURL, endpoint.QueryID, endpoint.OperationName)

	body := map[string]interface{}{
		"variables": variables,
		"queryId":   endpoint.QueryID,
	}
	if endpoint.HasFeatures {
		body["features"] = DefaultFeatures
	}
	if endpoint.HasFieldToggles {
		body["fieldToggles"] = DefaultFieldToggles
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(string(bodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)
	return c.doRequest(req)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", c.credentials.BearerToken)
	req.Header.Set("Cookie", c.credentials.Cookies)
	req.Header.Set("X-Csrf-Token", c.credentials.CSRFToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Client-Language", "en")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse rate limit headers
	c.LastRateLimit = ParseRateLimit(resp.Header)

	if resp.StatusCode == 429 {
		if c.LastRateLimit != nil {
			return nil, fmt.Errorf("rate limited — resets at %s", c.LastRateLimit.Reset.Format("15:04:05"))
		}
		return nil, fmt.Errorf("rate limited — try again later")
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, fmt.Errorf("authentication failed (HTTP %d) — try: x-cli auth login", resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, truncate(string(data), 200))
	}

	return data, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
