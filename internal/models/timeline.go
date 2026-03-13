package models

import (
	"github.com/tidwall/gjson"
)

// TimelinePage holds a page of tweets with pagination cursors.
type TimelinePage struct {
	Tweets     []*Tweet
	NextCursor string
	PrevCursor string
}

// ParseTimeline extracts tweets and cursors from a timeline instructions array.
// instructionsPath should point to the instructions array in the response JSON.
func ParseTimeline(instructions gjson.Result) TimelinePage {
	page := TimelinePage{}

	instructions.ForEach(func(_, instr gjson.Result) bool {
		instrType := instr.Get("type").String()

		var entries gjson.Result
		switch instrType {
		case "TimelineAddEntries":
			entries = instr.Get("entries")
		case "TimelineAddToModule":
			// Module items (e.g., conversation threads)
			instr.Get("moduleItems").ForEach(func(_, item gjson.Result) bool {
				tweetResult := item.Get("item.itemContent.tweet_results.result")
				if tweetResult.Exists() {
					if t := ParseTweet(tweetResult); t != nil {
						page.Tweets = append(page.Tweets, t)
					}
				}
				return true
			})
			return true
		default:
			return true
		}

		entries.ForEach(func(_, entry gjson.Result) bool {
			content := entry.Get("content")
			entryType := content.Get("entryType").String()

			switch entryType {
			case "TimelineTimelineItem":
				tweetResult := content.Get("itemContent.tweet_results.result")
				if tweetResult.Exists() {
					if t := ParseTweet(tweetResult); t != nil {
						page.Tweets = append(page.Tweets, t)
					}
				}
			case "TimelineTimelineModule":
				// Module with items (e.g., conversation threads)
				content.Get("items").ForEach(func(_, item gjson.Result) bool {
					tweetResult := item.Get("item.itemContent.tweet_results.result")
					if tweetResult.Exists() {
						if t := ParseTweet(tweetResult); t != nil {
							page.Tweets = append(page.Tweets, t)
						}
					}
					return true
				})
			case "TimelineTimelineCursor":
				cursorType := content.Get("cursorType").String()
				cursorValue := content.Get("value").String()
				switch cursorType {
				case "Bottom":
					page.NextCursor = cursorValue
				case "Top":
					page.PrevCursor = cursorValue
				}
			}
			return true
		})
		return true
	})

	return page
}
