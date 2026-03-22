package scraper

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/openai/openai-go/v2"
	"gopkg.in/yaml.v3"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
	"github.com/thrgamon/coffeeroasters/internal/normalise"
)

// RunOptions configures a scrape run.
type RunOptions struct {
	ConfigPath  string
	Slugs       []string // empty = all active roasters
	DryRun      bool
	Concurrency int
	Verbose     bool
}

// RunResult summarises a single roaster scrape.
type RunResult struct {
	Slug     string
	Coffees  int
	Duration time.Duration
	Err      error
}

// Runner orchestrates scraping across roasters.
type Runner struct {
	queries   *db.Queries
	pool      *pgxpool.Pool
	extractor *Extractor
	limiter   *RateLimiter
	robots    *RobotsCache
	client    *http.Client
}

// NewRunner creates a Runner. If pool is nil (dry-run mode), DB operations
// are skipped.
func NewRunner(pool *pgxpool.Pool, openaiClient *openai.Client) *Runner {
	var queries *db.Queries
	if pool != nil {
		queries = db.New(pool)
	}
	return &Runner{
		queries:   queries,
		pool:      pool,
		extractor: NewExtractor(openaiClient),
		limiter:   NewRateLimiter(),
		robots:    NewRobotsCache(),
		client:    DefaultHTTPClient(),
	}
}

// Run loads roaster configs and scrapes them concurrently.
func (r *Runner) Run(ctx context.Context, opts RunOptions) []RunResult {
	configs, err := loadRoasterConfigs(opts.ConfigPath)
	if err != nil {
		slog.Error("load configs", "error", err)
		return []RunResult{{Err: err}}
	}

	// Filter by slugs if specified
	if len(opts.Slugs) > 0 {
		slugSet := make(map[string]bool, len(opts.Slugs))
		for _, s := range opts.Slugs {
			slugSet[s] = true
		}
		var filtered []domain.RoasterConfig
		for _, c := range configs {
			if slugSet[c.Slug] {
				filtered = append(filtered, c)
			}
		}
		configs = filtered
	}

	// Only active roasters unless filtering by slug
	if len(opts.Slugs) == 0 {
		var active []domain.RoasterConfig
		for _, c := range configs {
			if c.Active {
				active = append(active, c)
			}
		}
		configs = active
	}

	if len(configs) == 0 {
		slog.Warn("no roasters to scrape")
		return nil
	}

	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = 3
	}
	if concurrency > 10 {
		concurrency = 10
	}

	sem := make(chan struct{}, concurrency)
	var mu sync.Mutex
	var results []RunResult

	var wg sync.WaitGroup
	for _, cfg := range configs {
		wg.Add(1)
		go func(cfg domain.RoasterConfig) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := r.scrapeOne(ctx, cfg, opts)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(cfg)
	}
	wg.Wait()

	return results
}

func (r *Runner) scrapeOne(ctx context.Context, cfg domain.RoasterConfig, opts RunOptions) RunResult {
	start := time.Now()
	logger := slog.With("roaster", cfg.Slug)

	// Check robots.txt
	allowed, err := r.robots.IsAllowed(ctx, cfg.ShopURL)
	if err != nil {
		logger.Warn("robots.txt check failed, proceeding", "error", err)
	}
	if !allowed {
		logger.Info("blocked by robots.txt, skipping")
		return RunResult{
			Slug:     cfg.Slug,
			Duration: time.Since(start),
			Err:      fmt.Errorf("blocked by robots.txt"),
		}
	}

	// Fetch and extract
	var coffees []RawCoffee
	switch cfg.FetchMethod {
	case domain.FetchShopifyJSON:
		coffees, err = FetchShopify(ctx, cfg, r.client, r.extractor, r.limiter)
	case domain.FetchHTML:
		coffees, err = FetchHTMLPage(ctx, cfg, r.client, r.extractor, r.limiter)
	case domain.FetchHTMLDetail:
		coffees, err = FetchHTMLDetail(ctx, cfg, r.client, r.extractor, r.limiter)
	default:
		err = fmt.Errorf("unknown fetch method: %s", cfg.FetchMethod)
	}

	if err != nil {
		logger.Error("scrape failed", "error", err)
		return RunResult{
			Slug:     cfg.Slug,
			Duration: time.Since(start),
			Err:      err,
		}
	}

	logger.Info("scraped", "coffees", len(coffees))

	if opts.DryRun {
		for _, c := range coffees {
			countryCode, regionName := normalise.NormaliseOrigin(c.OriginRaw, c.RegionRaw)
			fmt.Printf("  %s | %s | %s (country=%s region=%s) | %s | %s | %s | producer=%s\n",
				c.Name, c.OriginRaw, c.RegionRaw, countryCode, regionName,
				c.ProcessRaw, c.RoastRaw, c.PriceRaw, c.ProducerRaw)
		}
		return RunResult{
			Slug:     cfg.Slug,
			Coffees:  len(coffees),
			Duration: time.Since(start),
		}
	}

	// Upsert to DB
	if r.queries != nil {
		if err := r.upsertRoasterAndCoffees(ctx, cfg, coffees); err != nil {
			logger.Error("db upsert failed", "error", err)
			return RunResult{
				Slug:     cfg.Slug,
				Coffees:  len(coffees),
				Duration: time.Since(start),
				Err:      err,
			}
		}
	}

	return RunResult{
		Slug:     cfg.Slug,
		Coffees:  len(coffees),
		Duration: time.Since(start),
	}
}

func (r *Runner) upsertRoasterAndCoffees(ctx context.Context, cfg domain.RoasterConfig, coffees []RawCoffee) error {
	start := time.Now()

	// Upsert roaster
	roasterID, err := r.queries.UpsertRoaster(ctx, db.UpsertRoasterParams{
		Slug:    cfg.Slug,
		Name:    cfg.Name,
		Website: cfg.Website,
		State:   textVal(cfg.State),
	})
	if err != nil {
		return fmt.Errorf("upsert roaster %s: %w", cfg.Slug, err)
	}

	// Record scrape run
	runID, err := r.queries.InsertScrapeRun(ctx, roasterID)
	if err != nil {
		return fmt.Errorf("insert scrape run: %w", err)
	}

	// Upsert each coffee
	var added, updated int
	for _, raw := range coffees {
		sanitizeRawCoffee(&raw)

		process := normalise.NormaliseProcess(raw.ProcessRaw)
		roastLevel := normalise.NormaliseRoastLevel(raw.RoastRaw)
		tastingNotes := normalise.NormaliseTastingNotes(raw.TastingNotes)
		priceCents, _ := normalise.NormalisePriceAUD(raw.PriceRaw)
		weightGrams, _ := normalise.NormaliseWeightGrams(raw.WeightRaw)

		// Per-100g pricing: use pre-computed values from Shopify variants,
		// or compute from single price/weight for HTML paths.
		per100gMin := raw.PricePer100gMin
		per100gMax := raw.PricePer100gMax
		if per100gMin == 0 && priceCents > 0 && weightGrams > 0 {
			p := normalise.PricePer100g(priceCents, weightGrams)
			per100gMin = p
			per100gMax = p
		}

		// Origin normalisation (skip for blends)
		var countryCode, regionName string
		if !raw.IsBlend {
			countryCode, regionName = normalise.NormaliseOrigin(raw.OriginRaw, raw.RegionRaw)
		}

		// Variety normalisation
		variety, species := normalise.NormaliseVariety(raw.VarietyRaw)

		var regionID pgtype.Int4
		if countryCode != "" && regionName != "" {
			rid, err := r.queries.GetOrCreateRegion(ctx, db.GetOrCreateRegionParams{
				CountryCode: countryCode,
				Name:        regionName,
			})
			if err != nil {
				slog.Warn("get-or-create region failed", "country", countryCode, "region", regionName, "error", err)
			} else {
				regionID = int4Val(rid)
			}
		}

		var producerID pgtype.Int4
		if raw.ProducerRaw != "" && countryCode != "" {
			pid, err := r.queries.GetOrCreateProducer(ctx, db.GetOrCreateProducerParams{
				Name:        raw.ProducerRaw,
				CountryCode: textVal(countryCode),
				RegionID:    regionID,
			})
			if err != nil {
				slog.Warn("get-or-create producer failed", "producer", raw.ProducerRaw, "error", err)
			} else {
				producerID = int4Val(pid)
			}
		}

		upsertResult, err := r.queries.UpsertCoffee(ctx, db.UpsertCoffeeParams{
			RoasterID:       roasterID,
			Name:            raw.Name,
			ProductUrl:      textVal(raw.ProductURL),
			ImageUrl:        textVal(raw.ImageURL),
			OriginRaw:       textVal(raw.OriginRaw),
			RegionRaw:       textVal(raw.RegionRaw),
			VarietyRaw:      textVal(raw.VarietyRaw),
			ProcessRaw:      textVal(raw.ProcessRaw),
			RoastRaw:        textVal(raw.RoastRaw),
			TastingNotesRaw: textVal(raw.TastingNotes),
			PriceRaw:        textVal(raw.PriceRaw),
			WeightRaw:       textVal(raw.WeightRaw),
			Currency:        raw.Currency,
			InStock:         raw.InStock,
			Process:         textVal(process),
			RoastLevel:      textVal(roastLevel),
			TastingNotes:    tastingNotes,
			PriceCents:      int4Val(int32(priceCents)),
			WeightGrams:     int4Val(int32(weightGrams)),
			CountryCode:     textVal(countryCode),
			RegionID:        regionID,
			ProducerID:      producerID,
			ProducerRaw:     textVal(raw.ProducerRaw),
			Variety:         textVal(variety),
			Species:         textVal(species),
			PricePer100gMin: int4Val(int32(per100gMin)),
			PricePer100gMax: int4Val(int32(per100gMax)),
			IsBlend:         raw.IsBlend,
			Description:     textVal(raw.Description),
		})
		if err != nil {
			slog.Warn("upsert coffee failed", "name", raw.Name, "error", err)
			continue
		}

		coffeeID := upsertResult.ID

		// Handle blend components: clean re-insert each scrape
		if raw.IsBlend && len(raw.BlendComponents) > 0 {
			if err := r.queries.DeleteBlendComponents(ctx, int32(coffeeID)); err != nil {
				slog.Warn("delete blend components failed", "coffee", raw.Name, "error", err)
			}

			for _, comp := range raw.BlendComponents {
				compCountry, compRegion := normalise.NormaliseOrigin(comp.Origin, comp.Region)

				var compRegionID pgtype.Int4
				if compCountry != "" && compRegion != "" {
					rid, err := r.queries.GetOrCreateRegion(ctx, db.GetOrCreateRegionParams{
						CountryCode: compCountry,
						Name:        compRegion,
					})
					if err != nil {
						slog.Warn("get-or-create blend component region failed", "country", compCountry, "region", compRegion, "error", err)
					} else {
						compRegionID = int4Val(rid)
					}
				}

				if err := r.queries.InsertBlendComponent(ctx, db.InsertBlendComponentParams{
					CoffeeID:    int32(coffeeID),
					CountryCode: textVal(compCountry),
					RegionID:    compRegionID,
					Variety:     textVal(comp.Variety),
					Percentage:  int4Val(int32(comp.Percentage)),
				}); err != nil {
					slog.Warn("insert blend component failed", "coffee", raw.Name, "error", err)
				}
			}
		}

		if upsertResult.IsNew {
			added++
		} else {
			updated++
		}
	}

	// Update scrape run
	durationMs := int32(time.Since(start).Milliseconds())
	err = r.queries.CompleteScrapeRun(ctx, db.CompleteScrapeRunParams{
		ID:             runID,
		Status:         "success",
		CoffeesFound:   int4Val(int32(len(coffees))),
		CoffeesAdded:   int4Val(int32(added)),
		CoffeesUpdated: int4Val(int32(updated)),
		DurationMs:     int4Val(durationMs),
	})
	if err != nil {
		return fmt.Errorf("complete scrape run: %w", err)
	}

	return nil
}

func loadRoasterConfigs(path string) ([]domain.RoasterConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var f domain.RoastersFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	return f.Roasters, nil
}

// stripNullBytes removes null bytes that OpenAI structured output occasionally
// embeds in strings. Postgres rejects these with "invalid byte sequence for
// encoding UTF8: 0x00".
func stripNullBytes(s string) string {
	return strings.ReplaceAll(s, "\x00", "")
}

// sanitizeRawCoffee strips null bytes from all string fields on a RawCoffee,
// including nested RawBlendComponent strings.
func sanitizeRawCoffee(raw *RawCoffee) {
	raw.Name = stripNullBytes(raw.Name)
	raw.ProductURL = stripNullBytes(raw.ProductURL)
	raw.ImageURL = stripNullBytes(raw.ImageURL)
	raw.OriginRaw = stripNullBytes(raw.OriginRaw)
	raw.RegionRaw = stripNullBytes(raw.RegionRaw)
	raw.VarietyRaw = stripNullBytes(raw.VarietyRaw)
	raw.ProducerRaw = stripNullBytes(raw.ProducerRaw)
	raw.ProcessRaw = stripNullBytes(raw.ProcessRaw)
	raw.RoastRaw = stripNullBytes(raw.RoastRaw)
	raw.TastingNotes = stripNullBytes(raw.TastingNotes)
	raw.Description = stripNullBytes(raw.Description)
	raw.PriceRaw = stripNullBytes(raw.PriceRaw)
	raw.Currency = stripNullBytes(raw.Currency)
	raw.WeightRaw = stripNullBytes(raw.WeightRaw)

	for i := range raw.BlendComponents {
		raw.BlendComponents[i].Origin = stripNullBytes(raw.BlendComponents[i].Origin)
		raw.BlendComponents[i].Region = stripNullBytes(raw.BlendComponents[i].Region)
		raw.BlendComponents[i].Variety = stripNullBytes(raw.BlendComponents[i].Variety)
	}
}

func textVal(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func int4Val(v int32) pgtype.Int4 {
	if v == 0 {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: v, Valid: true}
}
