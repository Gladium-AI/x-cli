package api

// EndpointDef describes a GraphQL endpoint's operation details.
type EndpointDef struct {
	// QueryID is the hash-like identifier in the URL path. Rotates on X deploys.
	QueryID       string
	OperationName string
	Method        string // "GET" or "POST"
	HasFeatures   bool
	HasFieldToggles bool
}

// Endpoints maps operation names to their definitions.
// QueryIDs must be updated when X rotates them (on deploys).
// These can be extracted from X's JS bundles or captured from network traffic.
// Query IDs scraped from X's production JS bundle (main.*.js).
// Last updated: 2026-03-13
// These rotate on every X deploy — run `x-cli update-ids` or re-scrape manually.
var Endpoints = map[string]EndpointDef{
	"HomeTimeline": {
		QueryID:         "-HtXlyhboD0-JLXJ-xo9Vg",
		OperationName:   "HomeTimeline",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"UserTweets": {
		QueryID:         "Y59DTUMfcKmUAATiT2SlTw",
		OperationName:   "UserTweets",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"UserTweetsAndReplies": {
		QueryID:         "peEfJv5QXlXXlCECOPtHOQ",
		OperationName:   "UserTweetsAndReplies",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"TweetDetail": {
		QueryID:         "9rs110LSoPARDs61WOBZ7A",
		OperationName:   "TweetDetail",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"SearchTimeline": {
		QueryID:         "oKkjeoNFNQN7IeK7AHYc0A",
		OperationName:   "SearchTimeline",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"UserByScreenName": {
		QueryID:         "pLsOiyHJ1eFwPJlNmLp4Bg",
		OperationName:   "UserByScreenName",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"UsersByRestIds": {
		QueryID:         "8OKmcyotfczJb44QTTu5tQ",
		OperationName:   "UsersByRestIds",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"Followers": {
		QueryID:         "xBB-_3k-LNxWg8TFpuQiWQ",
		OperationName:   "Followers",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"Following": {
		QueryID:         "OEx3R66nP411LbwQ0xgAIg",
		OperationName:   "Following",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
	"BlueVerifiedFollowers": {
		QueryID:         "lsph9HGDm9-osG2BJO8RFg",
		OperationName:   "BlueVerifiedFollowers",
		Method:          "GET",
		HasFeatures:     true,
		HasFieldToggles: true,
	},
}
