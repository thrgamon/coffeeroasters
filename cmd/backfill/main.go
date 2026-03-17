// Command backfill populates country_code and region_id for existing coffees
// that have origin_raw set but no normalised country_code.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thrgamon/coffeeroasters/internal/config"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/normalise"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)

	rows, err := queries.ListCoffeesNeedingBackfill(ctx)
	if err != nil {
		slog.Error("list coffees needing backfill", "error", err)
		os.Exit(1)
	}

	slog.Info("coffees needing backfill", "count", len(rows))

	var updated, skipped int
	for _, row := range rows {
		countryCode, regionName := normalise.NormaliseOrigin(row.OriginRaw.String, row.RegionRaw.String)
		if countryCode == "" {
			skipped++
			continue
		}

		var regionID *int32
		if regionName != "" {
			rid, err := queries.GetOrCreateRegion(ctx, db.GetOrCreateRegionParams{
				CountryCode: countryCode,
				Name:        regionName,
			})
			if err != nil {
				slog.Warn("get-or-create region", "country", countryCode, "region", regionName, "error", err)
			} else {
				regionID = &rid
			}
		}

		params := db.UpdateCoffeeOriginParams{
			ID:          row.ID,
			CountryCode: textVal(countryCode),
		}
		if regionID != nil {
			params.RegionID = int4Val(*regionID)
		}

		if err := queries.UpdateCoffeeOrigin(ctx, params); err != nil {
			slog.Warn("update coffee origin", "id", row.ID, "error", err)
			continue
		}

		updated++
	}

	fmt.Printf("Backfill complete: %d updated, %d skipped (no country match)\n", updated, skipped)
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
