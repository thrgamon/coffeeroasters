package scraper

import (
	"context"
	"log/slog"
	"time"
)

// configDrivenScraper implements Scraper using a roaster Config.
// It dispatches to the appropriate strategy (jsonld, selector, shopify).
type configDrivenScraper struct {
	cfg         *Config
	rateLimiter *RateLimiter
	robotsCache *RobotsCache
	logger      *slog.Logger
	client      interface{ // *http.Client — typed as interface to allow test injection
	}
}

func newConfigDrivenScraper(cfg *Config, rl *RateLimiter, rc *RobotsCache, logger *slog.Logger) Scraper {
	return &configDrivenScraper{
		cfg:         cfg,
		rateLimiter: rl,
		robotsCache: rc,
		logger:      logger.With("roaster", cfg.Slug, "strategy", cfg.Strategy),
	}
}

func (s *configDrivenScraper) Slug() string {
	return s.cfg.Slug
}

func (s *configDrivenScraper) Scrape(ctx context.Context) (ScrapeResult, error) {
	result := ScrapeResult{
		RoasterSlug: s.cfg.Slug,
	}

	switch s.cfg.Strategy {
	case "shopify":
		coffees, pages, err := s.scrapeShopify(ctx)
		result.Coffees = coffees
		result.PagesVisited = pages
		if err != nil {
			result.Errors = append(result.Errors, err)
		}

	case "jsonld":
		coffees, pages, err := s.scrapeJSONLD(ctx)
		result.Coffees = coffees
		result.PagesVisited = pages
		if err != nil {
			result.Errors = append(result.Errors, err)
		}

	case "selector":
		coffees, pages, err := s.scrapeSelectors(ctx)
		result.Coffees = coffees
		result.PagesVisited = pages
		if err != nil {
			result.Errors = append(result.Errors, err)
		}

	default:
		result.Errors = append(result.Errors, newUnsupportedStrategyError(s.cfg.Strategy))
	}

	// Stamp scrape time on all results
	now := time.Now()
	for i := range result.Coffees {
		result.Coffees[i].ScrapedAt = now
		if result.Coffees[i].Currency == "" {
			result.Coffees[i].Currency = "AUD"
		}
	}

	// Apply skip patterns
	if len(s.cfg.FieldHints.SkipPatterns) > 0 {
		result.Coffees = filterSkipped(result.Coffees, s.cfg.FieldHints.SkipPatterns)
	}

	return result, nil
}

// scrapeShopify uses the Shopify products.json endpoint.
// TODO: implement in Phase 2.
func (s *configDrivenScraper) scrapeShopify(ctx context.Context) ([]RawCoffee, int, error) {
	s.logger.Info("shopify strategy not yet implemented — skipping")
	return nil, 0, nil
}

// scrapeJSONLD parses Schema.org Product JSON-LD from HTML pages.
// TODO: implement in Phase 2.
func (s *configDrivenScraper) scrapeJSONLD(ctx context.Context) ([]RawCoffee, int, error) {
	s.logger.Info("jsonld strategy not yet implemented — skipping")
	return nil, 0, nil
}

// scrapeSelectors uses CSS selectors defined in the config.
// TODO: implement in Phase 2.
func (s *configDrivenScraper) scrapeSelectors(ctx context.Context) ([]RawCoffee, int, error) {
	s.logger.Info("selector strategy not yet implemented — skipping")
	return nil, 0, nil
}

type unsupportedStrategyError struct{ strategy string }

func newUnsupportedStrategyError(s string) error {
	return &unsupportedStrategyError{s}
}
func (e *unsupportedStrategyError) Error() string {
	return "unsupported scrape strategy: " + e.strategy
}

func filterSkipped(coffees []RawCoffee, patterns []string) []RawCoffee {
	out := coffees[:0]
	for _, c := range coffees {
		skip := false
		for _, p := range patterns {
			// Case-insensitive substring match
			if containsFold(c.Name, p) {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, c)
		}
	}
	return out
}

func containsFold(s, substr string) bool {
	return len(s) >= len(substr) &&
		func() bool {
			sl, subl := []rune(s), []rune(substr)
			for i := 0; i <= len(sl)-len(subl); i++ {
				match := true
				for j, r := range subl {
					if toLowerRune(sl[i+j]) != toLowerRune(r) {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
			return false
		}()
}

func toLowerRune(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
