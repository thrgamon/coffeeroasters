// Package scraper provides the core interfaces and types for the
// Coffeeroasters scraping framework.
//
// Each roaster is driven by a Config (loaded from roasters/configs/*.yaml)
// and run through the Runner, which handles robots.txt checking, per-domain
// rate limiting, and structured result collection.
package scraper

import (
	"context"
	"time"
)

// RawCoffee is the unprocessed data extracted from a single product listing.
// All fields are raw strings — normalisation happens in the normalise package.
type RawCoffee struct {
	// Identity
	Name       string
	ProductURL string
	ImageURL   string

	// Origin info (raw strings, normalised later)
	OriginRaw   string // e.g. "Ethiopia Yirgacheffe" or "Ethiopia, Colombia"
	RegionRaw   string // e.g. "Gedeo Zone"
	VarietyRaw  string // e.g. "Heirloom, 74158"
	ProducerRaw string // e.g. "Finca El Paraiso", "Dumerso Washing Station"

	// Process and roast (raw strings, normalised later)
	ProcessRaw   string // e.g. "Natural", "Washed", "Honey"
	RoastRaw     string // e.g. "Light", "Filter Roast", "Medium-Light"
	TastingNotes string // e.g. "Blueberry, dark chocolate, jasmine"

	// Pricing
	PriceRaw    string // e.g. "$32.00", "32", "32.00"
	Currency    string // e.g. "AUD" — default to AUD if empty
	WeightRaw   string // e.g. "250g", "1kg", "250"
	InStock     bool

	// Metadata
	ScrapedAt time.Time
}

// RawRoaster is the unprocessed data for a roaster discovered via the
// discovery crawler.
type RawRoaster struct {
	Name        string
	Website     string
	State       string // AU state: VIC, NSW, QLD, WA, SA, TAS, ACT, NT
	Description string
	SourceURL   string // the directory page where we found them
	DiscoveredAt time.Time
}

// ScrapeResult is returned by a single scraper run.
type ScrapeResult struct {
	RoasterSlug string
	Coffees     []RawCoffee
	Errors      []error
	Duration    time.Duration
	PagesVisited int
}

// Scraper is the interface every roaster-specific scraper must implement.
// The config-driven runner constructs scrapers from YAML configs; bespoke
// scrapers implement this interface directly for complex sites.
type Scraper interface {
	// Scrape fetches and extracts all coffee listings for the roaster.
	Scrape(ctx context.Context) (ScrapeResult, error)

	// Slug returns the unique identifier for the roaster (matches the YAML filename).
	Slug() string
}
