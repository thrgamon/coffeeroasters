package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const botUserAgent = "CoffeeroastersBot/1.0 (+https://coffeeroasters.app)"

// RobotsCache fetches and caches robots.txt files per domain.
// We check robots.txt once per domain per process lifetime and cache the
// result to avoid unnecessary requests.
type RobotsCache struct {
	mu    sync.RWMutex
	cache map[string]*robotsEntry
}

type robotsEntry struct {
	disallowed []string
	fetchedAt  time.Time
}

// NewRobotsCache creates a new RobotsCache.
func NewRobotsCache() *RobotsCache {
	return &RobotsCache{
		cache: make(map[string]*robotsEntry),
	}
}

// IsAllowed returns true if our bot user agent is allowed to fetch rawURL
// according to the domain's robots.txt. Returns true (with an error) if
// robots.txt cannot be fetched — we log but don't hard-fail.
func (rc *RobotsCache) IsAllowed(ctx context.Context, rawURL string) (bool, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return true, fmt.Errorf("parse URL %q: %w", rawURL, err)
	}

	domain := u.Hostname()
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, domain)

	entry, err := rc.fetchEntry(ctx, domain, robotsURL)
	if err != nil {
		// Can't fetch robots.txt — proceed but warn.
		return true, fmt.Errorf("fetch robots.txt for %s: %w", domain, err)
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	for _, disallowed := range entry.disallowed {
		if strings.HasPrefix(path, disallowed) {
			return false, nil
		}
	}

	return true, nil
}

func (rc *RobotsCache) fetchEntry(ctx context.Context, domain, robotsURL string) (*robotsEntry, error) {
	// Check cache first (read lock)
	rc.mu.RLock()
	entry, ok := rc.cache[domain]
	rc.mu.RUnlock()
	if ok {
		return entry, nil
	}

	// Fetch robots.txt
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", botUserAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	entry = &robotsEntry{fetchedAt: time.Now()}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		if err != nil {
			return nil, err
		}
		entry.disallowed = parseRobotsTxt(string(body))
	}
	// 404 or other non-200 → no restrictions, empty disallowed list

	rc.mu.Lock()
	rc.cache[domain] = entry
	rc.mu.Unlock()

	return entry, nil
}

// parseRobotsTxt is a minimal robots.txt parser. It extracts Disallow rules
// that apply to our bot or to all user agents (*).
//
// This is intentionally simple — for a full implementation consider the
// golang.org/x/net/webdav or a dedicated robots.txt library.
func parseRobotsTxt(content string) []string {
	var disallowed []string
	applicable := false

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)

		// Strip comments
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "user-agent:") {
			agent := strings.TrimSpace(line[len("user-agent:"):])
			applicable = agent == "*" ||
				strings.EqualFold(agent, "CoffeeroastersBot") ||
				strings.EqualFold(agent, "coffeeroastersbot/1.0")
			continue
		}

		if applicable && strings.HasPrefix(lower, "disallow:") {
			path := strings.TrimSpace(line[len("disallow:"):])
			if path != "" && path != "/" {
				disallowed = append(disallowed, path)
			} else if path == "/" {
				// Disallow everything
				disallowed = append(disallowed, "/")
			}
		}
	}

	return disallowed
}

// DefaultHTTPClient returns a pre-configured http.Client with our bot
// User-Agent, a sensible timeout, and no cookie jar (stateless scraping).
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			req.Header.Set("User-Agent", botUserAgent)
			return nil
		},
	}
}

// FetchHTML fetches a URL and returns the parsed HTML node tree.
func FetchHTML(ctx context.Context, client *http.Client, rawURL string) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", botUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", rawURL, resp.StatusCode)
	}

	doc, err := html.Parse(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("parse HTML from %s: %w", rawURL, err)
	}

	return doc, nil
}
