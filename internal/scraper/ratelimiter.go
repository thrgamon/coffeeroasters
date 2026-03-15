package scraper

import (
	"net/url"
	"sync"
	"time"
)

// RateLimiter enforces a per-domain minimum delay between HTTP requests.
// This is a core part of our ethical scraping policy — we never hammer
// a roaster's site, regardless of how fast our scrapers could go.
type RateLimiter struct {
	mu      sync.Mutex
	lastReq map[string]time.Time
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		lastReq: make(map[string]time.Time),
	}
}

// Wait blocks until the minimum delay has elapsed since the last request
// to the same domain as rawURL. minDelay is the per-config rate limit.
//
// If minDelay < 2 seconds it is clamped to 2 seconds as a safety floor.
func (r *RateLimiter) Wait(rawURL string, minDelay time.Duration) {
	if minDelay < 2*time.Second {
		minDelay = 2 * time.Second
	}

	domain := domainOf(rawURL)

	r.mu.Lock()
	last, ok := r.lastReq[domain]
	r.mu.Unlock()

	if ok {
		elapsed := time.Since(last)
		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	r.mu.Lock()
	r.lastReq[domain] = time.Now()
	r.mu.Unlock()
}

// Record records a request to rawURL without blocking. Used to account for
// requests made outside the rate limiter (e.g. robots.txt fetches).
func (r *RateLimiter) Record(rawURL string) {
	domain := domainOf(rawURL)
	r.mu.Lock()
	r.lastReq[domain] = time.Now()
	r.mu.Unlock()
}

func domainOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Hostname()
}
