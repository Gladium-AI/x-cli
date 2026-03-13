package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RateLimit holds parsed rate limit info from response headers.
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// ParseRateLimit extracts rate limit info from response headers.
func ParseRateLimit(h http.Header) *RateLimit {
	limit, _ := strconv.Atoi(h.Get("X-Rate-Limit-Limit"))
	remaining, _ := strconv.Atoi(h.Get("X-Rate-Limit-Remaining"))
	resetUnix, _ := strconv.ParseInt(h.Get("X-Rate-Limit-Reset"), 10, 64)

	if limit == 0 && remaining == 0 {
		return nil
	}

	return &RateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     time.Unix(resetUnix, 0),
	}
}

// WaitIfNeeded blocks if rate limit is exhausted, printing a warning.
func (rl *RateLimit) WaitIfNeeded() {
	if rl == nil || rl.Remaining > 0 {
		return
	}
	waitDuration := time.Until(rl.Reset)
	if waitDuration <= 0 {
		return
	}
	fmt.Printf("Rate limited. Waiting %s...\n", waitDuration.Round(time.Second))
	time.Sleep(waitDuration + time.Second)
}
