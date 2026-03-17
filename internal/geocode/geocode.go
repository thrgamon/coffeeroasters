package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/llm"
)

// RegionCoords holds the approximate latitude and longitude of a coffee-growing region.
type RegionCoords struct {
	Latitude  float64 `json:"latitude" jsonschema:"description=Approximate latitude of the centre of the growing region"`
	Longitude float64 `json:"longitude" jsonschema:"description=Approximate longitude of the centre of the growing region"`
}

// Geocoder uses GPT-4.1-mini with structured output to approximate the GPS
// coordinates of coffee-growing regions.
type Geocoder struct {
	client *openai.Client
}

// NewGeocoder creates a Geocoder from an OpenAI client.
func NewGeocoder(client *openai.Client) *Geocoder {
	return &Geocoder{client: client}
}

const geocodeSystemPrompt = `Return the approximate latitude and longitude of the given coffee-growing region. Use the centre of the growing area.`

// Geocode returns approximate GPS coordinates for a coffee-growing region.
func (g *Geocoder) Geocode(ctx context.Context, regionName, countryName string) (RegionCoords, error) {
	schema := llm.GenerateSchema(RegionCoords{})

	resp, err := g.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(geocodeSystemPrompt),
			openai.UserMessage(fmt.Sprintf("%s, %s", regionName, countryName)),
		},
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "region_coords",
					Strict: openai.Opt(true),
					Schema: schema,
				},
			},
		},
	})
	if err != nil {
		return RegionCoords{}, fmt.Errorf("openai geocode: %w", err)
	}

	if len(resp.Choices) == 0 {
		return RegionCoords{}, fmt.Errorf("openai returned no choices")
	}

	var coords RegionCoords
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &coords); err != nil {
		return RegionCoords{}, fmt.Errorf("parse geocode response: %w", err)
	}

	return coords, nil
}

// BackfillPending geocodes all regions that are missing coordinates and stores
// the results. Returns counts of geocoded and failed regions.
func (g *Geocoder) BackfillPending(ctx context.Context, queries *db.Queries) (geocoded, failed int) {
	rows, err := queries.ListRegionsNeedingGeocode(ctx)
	if err != nil {
		slog.Error("list regions needing geocode", "error", err)
		return 0, 0
	}

	if len(rows) == 0 {
		return 0, 0
	}

	slog.Info("regions needing geocode", "count", len(rows))

	for _, row := range rows {
		coords, err := g.Geocode(ctx, row.Name, row.CountryName)
		if err != nil {
			slog.Warn("geocode failed", "region", row.Name, "country", row.CountryName, "error", err)
			failed++
			continue
		}

		err = queries.UpdateRegionCoordinates(ctx, db.UpdateRegionCoordinatesParams{
			ID:        row.ID,
			Latitude:  pgtype.Float8{Float64: coords.Latitude, Valid: true},
			Longitude: pgtype.Float8{Float64: coords.Longitude, Valid: true},
		})
		if err != nil {
			slog.Warn("update coordinates", "region", row.Name, "error", err)
			failed++
			continue
		}

		slog.Info("geocoded region", "region", row.Name, "country", row.CountryName,
			"lat", fmt.Sprintf("%.4f", coords.Latitude), "lon", fmt.Sprintf("%.4f", coords.Longitude))
		geocoded++

		time.Sleep(500 * time.Millisecond)
	}

	return geocoded, failed
}
