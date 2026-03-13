package models

import (
	"strconv"

	"github.com/tidwall/gjson"
)

// User represents a parsed X user profile.
type User struct {
	ID              string
	ScreenName      string
	Name            string
	Description     string
	Location        string
	URL             string
	FollowersCount  int
	FriendsCount    int
	StatusesCount   int
	FavouritesCount int
	MediaCount      int
	ListedCount     int
	IsBlueVerified  bool
	IsVerified      bool
	IsProtected     bool
	ProfileImageURL string
	ProfileBannerURL string
	CreatedAt       string
	PinnedTweetIDs  []string
}

// ParseUser extracts a User from a gjson result pointing at a user result object.
// Expected path: data.user.result (for UserByScreenName)
func ParseUser(r gjson.Result) User {
	u := User{
		ID:               r.Get("rest_id").String(),
		Name:             r.Get("core.name").String(),
		ScreenName:       r.Get("core.screen_name").String(),
		CreatedAt:        r.Get("core.created_at").String(),
		IsBlueVerified:   r.Get("is_blue_verified").Bool(),
		IsVerified:       r.Get("verification.verified").Bool(),
		IsProtected:      r.Get("privacy.protected").Bool(),
		ProfileImageURL:  r.Get("avatar.image_url").String(),
		Description:      r.Get("legacy.description").String(),
		FollowersCount:   int(r.Get("legacy.followers_count").Int()),
		FriendsCount:     int(r.Get("legacy.friends_count").Int()),
		StatusesCount:    int(r.Get("legacy.statuses_count").Int()),
		FavouritesCount:  int(r.Get("legacy.favourites_count").Int()),
		MediaCount:       int(r.Get("legacy.media_count").Int()),
		ListedCount:      int(r.Get("legacy.listed_count").Int()),
		ProfileBannerURL: r.Get("legacy.profile_banner_url").String(),
		Location:         r.Get("location.location").String(),
	}

	// Extract expanded URL if available
	expandedURL := r.Get("legacy.entities.url.urls.0.expanded_url").String()
	if expandedURL != "" {
		u.URL = expandedURL
	}

	// Pinned tweets
	r.Get("legacy.pinned_tweet_ids_str").ForEach(func(_, v gjson.Result) bool {
		u.PinnedTweetIDs = append(u.PinnedTweetIDs, v.String())
		return true
	})

	return u
}

// UserSummary is a compact user representation embedded in tweet results.
type UserSummary struct {
	ID             string
	ScreenName     string
	Name           string
	IsBlueVerified bool
	ProfileImageURL string
}

// ParseUserSummary extracts a compact user from a tweet's core.user_results.result.
func ParseUserSummary(r gjson.Result) UserSummary {
	return UserSummary{
		ID:              r.Get("rest_id").String(),
		Name:            r.Get("core.name").String(),
		ScreenName:      r.Get("core.screen_name").String(),
		IsBlueVerified:  r.Get("is_blue_verified").Bool(),
		ProfileImageURL: r.Get("avatar.image_url").String(),
	}
}

// FormatCount formats a number into a human-readable string (e.g., 1.2K, 5.4M).
func FormatCount(n int) string {
	switch {
	case n >= 1_000_000_000:
		return strconv.FormatFloat(float64(n)/1_000_000_000, 'f', 1, 64) + "B"
	case n >= 1_000_000:
		return strconv.FormatFloat(float64(n)/1_000_000, 'f', 1, 64) + "M"
	case n >= 10_000:
		return strconv.FormatFloat(float64(n)/1_000, 'f', 1, 64) + "K"
	case n >= 1_000:
		return strconv.FormatFloat(float64(n)/1_000, 'f', 1, 64) + "K"
	default:
		return strconv.Itoa(n)
	}
}
