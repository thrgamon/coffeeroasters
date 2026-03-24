package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/openai/openai-go/v2"

	"github.com/thrgamon/coffeeroasters/internal/classify"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/embedding"
	"github.com/thrgamon/coffeeroasters/internal/geocode"
	"github.com/thrgamon/coffeeroasters/internal/scraper"
)

func main() {
	var (
		roaster     = flag.String("roaster", "", "comma-separated roaster slugs (empty = all active)")
		dryRun      = flag.Bool("dry-run", false, "extract and print without writing to DB")
		configPath  = flag.String("config", "roasters.yaml", "path to roasters config file")
		concurrency = flag.Int("concurrency", 3, "max concurrent scrapes")
		verbose     = flag.Bool("verbose", false, "enable debug logging")
	)
	flag.Parse()

	if *verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	openaiClient := openai.NewClient()

	var pool *pgxpool.Pool
	if !*dryRun {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://postgres:postgres@localhost:5432/coffeeroasters?sslmode=disable"
		}

		var err error
		pool, err = pgxpool.New(ctx, dbURL)
		if err != nil {
			log.Fatalf("connect database: %v", err)
		}
		defer pool.Close()

		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := pool.Ping(pingCtx); err != nil {
			cancel()
			log.Fatalf("ping database: %v", err)
		}
		cancel()

		// Verify schema is up to date before running queries that depend on
		// newer columns. Fails fast with a clear message if migrations are pending.
		if err := verifySchema(ctx, pool); err != nil {
			log.Fatalf("%v", err)
		}
	}

	runner := scraper.NewRunner(pool, &openaiClient)

	var slugs []string
	if *roaster != "" {
		slugs = strings.Split(*roaster, ",")
	}

	start := time.Now()
	results := runner.Run(ctx, scraper.RunOptions{
		ConfigPath:  *configPath,
		Slugs:       slugs,
		DryRun:      *dryRun,
		Concurrency: *concurrency,
		Verbose:     *verbose,
	})

	var total, success, failed int
	for _, r := range results {
		total++
		if r.Err != nil {
			failed++
			slog.Error("roaster failed", "slug", r.Slug, "error", r.Err, "duration", r.Duration)
		} else {
			success++
			slog.Info("roaster done", "slug", r.Slug, "coffees", r.Coffees, "duration", r.Duration)
		}
	}

	fmt.Printf("\nScrape complete: %d total, %d success, %d failed, %s elapsed\n",
		total, success, failed, time.Since(start).Round(time.Millisecond))

	// Geocode any new regions discovered during scraping
	if !*dryRun && pool != nil {
		queries := db.New(pool)
		geocoder := geocode.NewGeocoder(&openaiClient)
		geocoded, geoFailed := geocoder.BackfillPending(ctx, queries)
		if geocoded > 0 || geoFailed > 0 {
			fmt.Printf("Geocode: %d geocoded, %d failed\n", geocoded, geoFailed)
		}

		// Classify unrecognised varieties via LLM
		classifier := classify.NewClassifier(&openaiClient)
		classified, classifyFailed := classifier.BackfillUnclassified(ctx, queries)
		if classified > 0 || classifyFailed > 0 {
			fmt.Printf("Variety classify: %d classified, %d failed\n", classified, classifyFailed)
		}

		// Embed descriptions for similarity scoring
		embedder := embedding.NewEmbedder(&openaiClient)
		embedded, embedFailed := embedder.BackfillPending(ctx, queries)
		if embedded > 0 || embedFailed > 0 {
			fmt.Printf("Embedding: %d embedded, %d failed\n", embedded, embedFailed)
		}
	}

	if failed > 0 {
		os.Exit(1)
	}
}

// verifySchema checks that the coffees table has all columns the scraper
// depends on. This catches missing migrations before we hit cryptic
// "column does not exist" errors mid-scrape.
func verifySchema(ctx context.Context, pool *pgxpool.Pool) error {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := pool.Exec(checkCtx, "SELECT description, source_hash FROM coffees LIMIT 0")
	if err != nil {
		return fmt.Errorf("schema out of date (run migrations): %w", err)
	}
	return nil
}
