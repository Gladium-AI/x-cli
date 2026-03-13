package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Tweet represents a parsed tweet.
type Tweet struct {
	ID           string
	FullText     string
	Author       UserSummary
	CreatedAt    time.Time
	RetweetCount int
	LikeCount    int
	ReplyCount   int
	QuoteCount   int
	BookmarkCount int
	ViewCount    int
	IsRetweet    bool
	IsQuote      bool
	Language     string
	InReplyTo    string
	Media        []Media
	QuotedTweet  *Tweet
}

// Media represents a media attachment on a tweet.
type Media struct {
	Type     string // "photo", "video", "animated_gif"
	URL      string
	Width    int
	Height   int
	VideoURL string // Best quality video variant URL
}

// ParseTweet extracts a Tweet from a gjson result pointing at a tweet_results.result.
func ParseTweet(r gjson.Result) *Tweet {
	// Handle TweetWithVisibilityResults wrapper
	typename := r.Get("__typename").String()
	if typename == "TweetWithVisibilityResults" {
		r = r.Get("tweet")
	}
	if typename == "TweetTombstone" {
		return nil
	}

	legacy := r.Get("legacy")
	if !legacy.Exists() {
		return nil
	}

	t := &Tweet{
		ID:            r.Get("rest_id").String(),
		Author:        ParseUserSummary(r.Get("core.user_results.result")),
		RetweetCount:  int(legacy.Get("retweet_count").Int()),
		LikeCount:     int(legacy.Get("favorite_count").Int()),
		ReplyCount:    int(legacy.Get("reply_count").Int()),
		QuoteCount:    int(legacy.Get("quote_count").Int()),
		BookmarkCount: int(legacy.Get("bookmark_count").Int()),
		IsQuote:       legacy.Get("is_quote_status").Bool(),
		Language:      legacy.Get("lang").String(),
		InReplyTo:     legacy.Get("in_reply_to_status_id_str").String(),
	}

	// Text: prefer note_tweet (long-form) over legacy.full_text
	noteText := r.Get("note_tweet.note_tweet_results.result.text").String()
	if noteText != "" {
		t.FullText = noteText
	} else {
		t.FullText = legacy.Get("full_text").String()
	}

	// View count
	viewStr := r.Get("views.count").String()
	if viewStr != "" {
		t.ViewCount, _ = strconv.Atoi(viewStr)
	}

	// Created at
	t.CreatedAt = parseTwitterTime(legacy.Get("created_at").String())

	// Check if this is a retweet
	if legacy.Get("retweeted_status_result.result").Exists() {
		t.IsRetweet = true
		// For retweets, parse the inner tweet for display
		inner := ParseTweet(legacy.Get("retweeted_status_result.result"))
		if inner != nil {
			inner.IsRetweet = false
			// Keep the RT author info but show inner tweet content
			rtAuthor := t.Author
			*t = *inner
			t.IsRetweet = true
			t.Author = inner.Author
			_ = rtAuthor // RT author available if needed
		}
	}

	// Quoted tweet
	if r.Get("quoted_status_result.result").Exists() {
		t.QuotedTweet = ParseTweet(r.Get("quoted_status_result.result"))
	}

	// Media from extended_entities (preferred) or entities
	mediaPath := "extended_entities.media"
	if !legacy.Get(mediaPath).Exists() {
		mediaPath = "entities.media"
	}
	legacy.Get(mediaPath).ForEach(func(_, v gjson.Result) bool {
		m := Media{
			Type:   v.Get("type").String(),
			URL:    v.Get("media_url_https").String(),
			Width:  int(v.Get("original_info.width").Int()),
			Height: int(v.Get("original_info.height").Int()),
		}
		// Get best quality video URL
		if m.Type == "video" || m.Type == "animated_gif" {
			var bestBitrate int64
			v.Get("video_info.variants").ForEach(func(_, variant gjson.Result) bool {
				if variant.Get("content_type").String() == "video/mp4" {
					bitrate := variant.Get("bitrate").Int()
					if bitrate > bestBitrate {
						bestBitrate = bitrate
						m.VideoURL = variant.Get("url").String()
					}
				}
				return true
			})
		}
		t.Media = append(t.Media, m)
		return true
	})

	return t
}

// TimeAgo returns a human-readable relative time string.
func (t *Tweet) TimeAgo() string {
	d := time.Since(t.CreatedAt)
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return strconv.Itoa(int(d.Minutes())) + "m"
	case d < 24*time.Hour:
		return strconv.Itoa(int(d.Hours())) + "h"
	case d < 7*24*time.Hour:
		return strconv.Itoa(int(d.Hours()/24)) + "d"
	default:
		return t.CreatedAt.Format("Jan 2, 2006")
	}
}

// twitterTimeFormat is the format X uses for created_at fields.
const twitterTimeFormat = "Mon Jan 02 15:04:05 -0700 2006"

func parseTwitterTime(s string) time.Time {
	s = strings.TrimSpace(s)
	t, err := time.Parse(twitterTimeFormat, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
