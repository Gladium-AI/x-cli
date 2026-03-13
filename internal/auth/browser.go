package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
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
		// Hide automation indicators from the browser
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
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

	// Remove navigator.webdriver flag before any page loads
	if err := chromedp.Run(timeoutCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		// Inject script to run on every new document (including after navigations/redirects)
		_, err := page.AddScriptToEvaluateOnNewDocument(`
			Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
		`).Do(ctx)
		return err
	})); err != nil {
		return nil, fmt.Errorf("setup anti-detection: %w", err)
	}

	// Navigate to login page
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate("https://x.com/i/flow/login"),
	); err != nil {
		return nil, fmt.Errorf("navigate to login: %w", err)
	}

	// Poll for login completion using CDP cookies (not page JS).
	// This avoids "execution context destroyed" errors during Google OAuth redirects,
	// since network.GetCookies works at the browser/protocol level, not the page level.
	var cookies []*network.Cookie
	pollInterval := 2 * time.Second
	deadline := time.Now().Add(loginTimeout)

	for time.Now().Before(deadline) {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("login timed out or browser was closed")
		default:
		}

		// Check cookies via CDP protocol — works regardless of page navigation state
		err := chromedp.Run(timeoutCtx, chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().WithURLs([]string{"https://x.com"}).Do(ctx)
			return err
		}))
		if err != nil {
			// Browser might have been closed by user
			if strings.Contains(err.Error(), "not found") ||
				strings.Contains(err.Error(), "target closed") ||
				strings.Contains(err.Error(), "connection reset") {
				return nil, fmt.Errorf("browser was closed before login completed")
			}
			// Transient error (e.g., during navigation), keep polling
			time.Sleep(pollInterval)
			continue
		}

		// Check if auth cookies are present
		hasAuth := false
		hasCT0 := false
		for _, c := range cookies {
			if c.Name == "auth_token" {
				hasAuth = true
			}
			if c.Name == "ct0" && c.Value != "" {
				hasCT0 = true
			}
		}

		if hasAuth && hasCT0 {
			// Login successful — wait a moment for all cookies to settle
			time.Sleep(2 * time.Second)

			// Re-fetch cookies to get the final set
			_ = chromedp.Run(timeoutCtx, chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				cookies, err = network.GetCookies().WithURLs([]string{"https://x.com"}).Do(ctx)
				return err
			}))
			break
		}

		time.Sleep(pollInterval)
	}

	// Build credentials from cookies
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
		return nil, fmt.Errorf("login failed: no CSRF token (ct0 cookie) found — did you complete login?")
	}

	// Try to extract screen name from the page (best-effort, may fail after OAuth redirects)
	var screenName string
	_ = chromedp.Run(timeoutCtx,
		chromedp.Evaluate(`document.querySelector('[data-testid="AppTabBar_Profile_Link"]')?.getAttribute("href")?.replace("/","") || ""`, &screenName),
	)
	if screenName != "" {
		creds.ScreenName = "@" + screenName
	}

	fmt.Println("Login successful! Closing browser...")
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
