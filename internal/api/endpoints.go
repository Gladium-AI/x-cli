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
var Endpoints = map[string]EndpointDef{
	"HomeTimeline": {
		QueryID:       "HJFjzBgCs16TqxewQOeLNg",
		OperationName: "HomeTimeline",
		Method:        "GET",
		HasFeatures:   true,
	},
	"UserTweets": {
		QueryID:       "H8OOoI-5ZE4NxgRr8lfyWg",
		OperationName: "UserTweets",
		Method:        "GET",
		HasFeatures:   true,
		HasFieldToggles: true,
	},
	"UserTweetsAndReplies": {
		QueryID:       "Q6aAvPfbCgxslajGMJfQPQ",
		OperationName: "UserTweetsAndReplies",
		Method:        "GET",
		HasFeatures:   true,
		HasFieldToggles: true,
	},
	"TweetDetail": {
		QueryID:       "nBS-WpgA6ZG0CyNHD517JQ",
		OperationName: "TweetDetail",
		Method:        "GET",
		HasFeatures:   true,
		HasFieldToggles: true,
	},
	"SearchTimeline": {
		QueryID:       "MJpyQGqgklrVl_0yelcbEQ",
		OperationName: "SearchTimeline",
		Method:        "GET",
		HasFeatures:   true,
	},
	"UserByScreenName": {
		QueryID:       "xmU6X_CKVnQ5lSrCbAmJsg",
		OperationName: "UserByScreenName",
		Method:        "GET",
		HasFeatures:   true,
		HasFieldToggles: true,
	},
	"UsersByRestIds": {
		QueryID:       "OsonVVsOgn_9vNy0MxBp1g",
		OperationName: "UsersByRestIds",
		Method:        "GET",
		HasFeatures:   true,
	},
	"Followers": {
		QueryID:       "djdTXDIk2qhd4OStqlUFeQ",
		OperationName: "Followers",
		Method:        "GET",
		HasFeatures:   true,
	},
	"Following": {
		QueryID:       "iSicc7LrzWGBgDPL0tM_TQ",
		OperationName: "Following",
		Method:        "GET",
		HasFeatures:   true,
	},
	"BlueVerifiedFollowers": {
		QueryID:       "tD_BnJijVySVIBWXPo4Jrw",
		OperationName: "BlueVerifiedFollowers",
		Method:        "GET",
		HasFeatures:   true,
	},
}
