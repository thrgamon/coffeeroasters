// cmd/discover crawls Australian specialty coffee directories and community
// sources to find new roasters not yet in the system.
//
// Discovered roasters are written to stdout as YAML stubs (one per line)
// suitable for review and addition to roasters/configs/.
//
// Usage:
//
//	# Discover from all configured sources
//	go run ./cmd/discover
//
//	# Discover from a specific source
//	go run ./cmd/discover -source specialty-coffee-au
//
//	# Output as JSON instead of YAML
//	go run ./cmd/discover -format json
//
//	# Limit to a specific AU state
//	go run ./cmd/discover -state VIC
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// discoverySource represents a website or API to crawl for roaster listings.
type discoverySource struct {
	Name string
	URL  string
	// TODO: implement per-source crawlers in Phase 1
}

// knownSources is the list of Australian specialty coffee directories and
// community resources to crawl for new roasters.
//
// Add new sources here as they are identified.
var knownSources = []discoverySource{
	{
		Name: "specialty-coffee-au",
		URL:  "https://www.scaa.asn.au/find-a-roaster",
	},
	{
		Name: "good-food-guide-roasters",
		URL:  "https://www.goodfood.com.au/eat-out/coffee",
	},
	{
		Name: "openstreetmap-au",
		URL:  "https://overpass-api.de/api/interpreter?data=[out:json];area[\"ISO3166-1\"=\"AU\"]->.au;node[craft=coffee_roasters](area.au);out;",
	},
}

func main() {
	var (
		sourceFlag = flag.String("source", "", "Specific source name to crawl. Empty = all sources.")
		stateFlag  = flag.String("state", "", "Filter results to an AU state (VIC, NSW, QLD, WA, SA, TAS, ACT, NT).")
		formatFlag = flag.String("format", "yaml", "Output format: yaml or json.")
		dryRun     = flag.Bool("dry-run", false, "Print discovered roasters without any DB writes.")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = ctx       // used when crawlers are implemented
	_ = sourceFlag
	_ = stateFlag
	_ = formatFlag
	_ = dryRun

	logger.Info("discovery crawler not yet implemented",
		"sources_configured", len(knownSources),
		"hint", "implement per-source crawlers in Phase 1",
	)

	// TODO Phase 1: implement per-source crawlers
	// Each source will produce []scraper.RawRoaster which are:
	//   1. De-duplicated by domain name (fuzzy match on name as fallback)
	//   2. Compared against existing roasters in the DB
	//   3. New ones written to stdout as YAML config stubs for human review
}
