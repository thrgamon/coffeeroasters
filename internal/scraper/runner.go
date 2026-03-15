package scraper

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// RunOptions configures a scrape run.
type RunOptions struct {
	// ConfigDir is the directory containing roaster YAML configs.
	ConfigDir string

	// Slugs restricts the run to specific roasters. If empty, all active
	// roasters are scraped.
	Slugs []string

	// DryRun prints extracted data without writing to the database.
	DryRun bool

	// Concurrency is the number of roasters to scrape in parallel.
	// Defaults to 3. Capped at 10 to avoid hammering sites.
	Concurrency int

	// IncludeInactive includes roasters marked active: false.
	IncludeInactive bool
}

// RunResult aggregates results from a full scrape run across multiple roasters.
type RunResult struct {
	Results  []ScrapeResult
	Total    int
	Success  int
	Failed   int
	Duration time.Duration
}

// Runner orchestrates scraping across all configured roasters.
type Runner struct {
	rateLimiter *RateLimiter
	robotsCache *RobotsCache
	logger      *slog.Logger
}

// NewRunner creates a new Runner.
func NewRunner(logger *slog.Logger) *Runner {
	return &Runner{
		rateLimiter: NewRateLimiter(),
		robotsCache: NewRobotsCache(),
		logger:      logger,
	}
}

// Run executes a full scrape run according to opts.
func (r *Runner) Run(ctx context.Context, opts RunOptions) (RunResult, error) {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 3
	}
	if opts.Concurrency > 10 {
		opts.Concurrency = 10
	}

	configs, err := LoadAllConfigs(opts.ConfigDir, opts.IncludeInactive)
	if err != nil {
		return RunResult{}, fmt.Errorf("load configs: %w", err)
	}

	// Filter by requested slugs
	if len(opts.Slugs) > 0 {
		slugSet := make(map[string]bool, len(opts.Slugs))
		for _, s := range opts.Slugs {
			slugSet[s] = true
		}
		filtered := configs[:0]
		for _, c := range configs {
			if slugSet[c.Slug] {
				filtered = append(filtered, c)
			}
		}
		configs = filtered
	}

	if len(configs) == 0 {
		return RunResult{}, fmt.Errorf("no matching active roaster configs found")
	}

	start := time.Now()
	sem := make(chan struct{}, opts.Concurrency)
	resultsCh := make(chan ScrapeResult, len(configs))

	var wg sync.WaitGroup
	for _, cfg := range configs {
		wg.Add(1)
		go func(cfg *Config) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := r.scrapeOne(ctx, cfg, opts.DryRun)
			resultsCh <- result
		}(cfg)
	}

	wg.Wait()
	close(resultsCh)

	var run RunResult
	run.Duration = time.Since(start)
	for res := range resultsCh {
		run.Results = append(run.Results, res)
		run.Total++
		if len(res.Errors) == 0 {
			run.Success++
		} else {
			run.Failed++
		}
	}

	return run, nil
}

func (r *Runner) scrapeOne(ctx context.Context, cfg *Config, dryRun bool) ScrapeResult {
	log := r.logger.With("roaster", cfg.Slug)
	start := time.Now()

	// Check robots.txt
	allowed, err := r.robotsCache.IsAllowed(ctx, cfg.ShopURL)
	if err != nil {
		log.Warn("robots.txt check failed, proceeding cautiously", "error", err)
	}
	if !allowed {
		return ScrapeResult{
			RoasterSlug: cfg.Slug,
			Errors:      []error{fmt.Errorf("robots.txt disallows scraping %s", cfg.ShopURL)},
			Duration:    time.Since(start),
		}
	}

	// Respect rate limit
	r.rateLimiter.Wait(cfg.ShopURL, time.Duration(cfg.RateLimitSeconds)*time.Second)

	scraper := newConfigDrivenScraper(cfg, r.rateLimiter, r.robotsCache, r.logger)
	result, err := scraper.Scrape(ctx)
	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	result.Duration = time.Since(start)

	if dryRun {
		log.Info("dry run result",
			"coffees_found", len(result.Coffees),
			"pages", result.PagesVisited,
			"duration", result.Duration.Round(time.Millisecond),
		)
		for i, c := range result.Coffees {
			log.Info("  coffee",
				"i", i+1,
				"name", c.Name,
				"origin", c.OriginRaw,
				"process", c.ProcessRaw,
				"price", c.PriceRaw,
				"in_stock", c.InStock,
			)
		}
	}

	return result
}
