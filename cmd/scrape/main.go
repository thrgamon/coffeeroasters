// cmd/scrape runs the coffee scraper for one or all configured roasters.
//
// Usage:
//
//	# Scrape all active roasters and write to DB
//	go run ./cmd/scrape
//
//	# Dry run a single roaster (print results, no DB writes)
//	go run ./cmd/scrape -roaster seven-seeds -dry-run
//
//	# Scrape multiple specific roasters
//	go run ./cmd/scrape -roaster seven-seeds,market-lane
//
//	# Include inactive roasters (e.g. for testing a WIP config)
//	go run ./cmd/scrape -roaster my-new-roaster -include-inactive -dry-run
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/thrgamon/coffeeroasters/internal/scraper"
)

func main() {
	var (
		roasterFlag      = flag.String("roaster", "", "Comma-separated roaster slugs to scrape. Empty = all active.")
		dryRun           = flag.Bool("dry-run", false, "Print results without writing to DB.")
		configDir        = flag.String("config-dir", "roasters/configs", "Directory containing roaster YAML configs.")
		concurrency      = flag.Int("concurrency", 3, "Number of roasters to scrape in parallel (max 10).")
		includeInactive  = flag.Bool("include-inactive", false, "Include roasters with active: false.")
		verbose          = flag.Bool("verbose", false, "Enable debug logging.")
	)
	flag.Parse()

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var slugs []string
	if *roasterFlag != "" {
		for _, s := range strings.Split(*roasterFlag, ",") {
			if trimmed := strings.TrimSpace(s); trimmed != "" {
				slugs = append(slugs, trimmed)
			}
		}
	}

	opts := scraper.RunOptions{
		ConfigDir:       *configDir,
		Slugs:           slugs,
		DryRun:          *dryRun,
		Concurrency:     *concurrency,
		IncludeInactive: *includeInactive,
	}

	runner := scraper.NewRunner(logger)
	result, err := runner.Run(ctx, opts)
	if err != nil {
		logger.Error("scrape run failed", "error", err)
		os.Exit(1)
	}

	logger.Info("scrape run complete",
		"total", result.Total,
		"success", result.Success,
		"failed", result.Failed,
		"duration", result.Duration,
	)

	if result.Failed > 0 {
		os.Exit(1)
	}
}
