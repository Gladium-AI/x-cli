package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Well-known bearer token used by X's web app
const defaultBearerToken = "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

const loginTimeout = 5 * time.Minute

// BrowserLogin opens a visible Chrome window for the user to log in to X,
// then extracts cookies, CSRF token, and basic user info.
func BrowserLogin(ctx context.Context) (*Credentials, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(500, 700),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(taskCtx, loginTimeout)
	defer timeoutCancel()

	fmt.Println("Opening browser for X login...")
	fmt.Println("Please log in to your X account. The window will close automatically.")

	// Navigate to login page
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate("https://x.com/i/flow/login"),
	); err != nil {
		return nil, fmt.Errorf("navigate to login: %w", err)
	}

	// Poll until we detect successful login (URL changes to /home or auth cookies appear)
	if err := chromedp.Run(timeoutCtx,
		chromedp.Poll(`window.location.pathname === "/home" || document.cookie.includes("auth_token")`, nil,
			chromedp.WithPollingInterval(500*time.Millisecond),
			chromedp.WithPollingTimeout(loginTimeout),
		),
	); err != nil {
		return nil, fmt.Errorf("waiting for login (timed out or browser closed): %w", err)
	}

	// Brief pause to let all cookies settle
	if err := chromedp.Run(timeoutCtx, chromedp.Sleep(2*time.Second)); err != nil {
		return nil, fmt.Errorf("post-login wait: %w", err)
	}

	// Extract cookies
	var cookies []*network.Cookie
	if err := chromedp.Run(timeoutCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().WithURLs([]string{"https://x.com"}).Do(ctx)
		return err
	})); err != nil {
		return nil, fmt.Errorf("extract cookies: %w", err)
	}

	// Build credential from cookies
	creds := &Credentials{
		BearerToken: defaultBearerToken,
		CreatedAt:   time.Now(),
	}

	var cookieParts []string
	for _, c := range cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
		switch c.Name {
		case "ct0":
			creds.CSRFToken = c.Value
		case "twid":
			// twid cookie is "u%3D<user_id>"
			creds.UserID = strings.TrimPrefix(c.Value, "u%3D")
		}
	}
	creds.Cookies = strings.Join(cookieParts, "; ")

	if creds.CSRFToken == "" {
		return nil, fmt.Errorf("login failed: no CSRF token (ct0 cookie) found")
	}

	// Try to extract screen name from the page
	var screenName string
	_ = chromedp.Run(timeoutCtx,
		chromedp.Evaluate(`document.querySelector('[data-testid="AppTabBar_Profile_Link"]')?.getAttribute("href")?.replace("/","")`, &screenName),
	)
	if screenName != "" {
		creds.ScreenName = "@" + screenName
	}

	return creds, nil
}

// CookiesToHTTP converts the stored cookie string into []*http.Cookie for use with http.Client.
func CookiesToHTTP(cookieStr string) []*http.Cookie {
	var cookies []*http.Cookie
	for _, part := range strings.Split(cookieStr, "; ") {
		parts := strings.SplitN(part, "=", 2)
		if len(parts) == 2 {
			cookies = append(cookies, &http.Cookie{Name: parts[0], Value: parts[1]})
		}
	}
	return cookies
}
