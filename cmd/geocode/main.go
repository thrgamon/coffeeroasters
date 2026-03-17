package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/openai/openai-go/v2"

	"github.com/thrgamon/coffeeroasters/internal/config"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/geocode"
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

	openaiClient := openai.NewClient()
	geocoder := geocode.NewGeocoder(&openaiClient)

	geocoded, failed := geocoder.BackfillPending(ctx, queries)
	fmt.Printf("Geocode complete: %d geocoded, %d failed\n", geocoded, failed)
}
