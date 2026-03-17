package classify

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

// ClassifiedVariety holds the structured output from the LLM classifier.
type ClassifiedVariety struct {
	Variety string `json:"variety" jsonschema:"description=Canonical lowercase variety name (e.g. bourbon, typica, gesha, sl28, heirloom). Empty string if unrecognisable."`
	Species string `json:"species" jsonschema:"description=Coffee species: arabica, robusta, liberica, or empty string if unknown."`
}

// Classifier uses GPT-4.1-mini with structured outputs to classify
// unrecognised coffee variety strings.
type Classifier struct {
	client *openai.Client
}

// NewClassifier creates a Classifier from an OpenAI client.
func NewClassifier(client *openai.Client) *Classifier {
	return &Classifier{client: client}
}

const classifySystemPrompt = `Given a raw coffee variety string, return the canonical variety name and species.
Use lowercase. Return empty strings if the input is not a recognisable coffee variety.

Common varieties: bourbon, typica, caturra, catuai, gesha, sl28, sl34, pacamara, maragogipe, heirloom, castillo, colombia, java, ruiru11, batian, catimor, marsellesa, parainema, obata, mundo-novo, yellow-bourbon, pink-bourbon, red-bourbon, tabi, sidra, wush-wush, 74110, 74112, 74158, pacas, maracaturra.
Species is one of: arabica, robusta, liberica, or empty string.`

// ClassifyVariety uses the LLM to classify a raw variety string into a
// canonical variety and species.
func (c *Classifier) ClassifyVariety(ctx context.Context, varietyRaw string) (ClassifiedVariety, error) {
	schema := llm.GenerateSchema(ClassifiedVariety{})

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(classifySystemPrompt),
			openai.UserMessage(varietyRaw),
		},
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "classified_variety",
					Strict: openai.Opt(true),
					Schema: schema,
				},
			},
		},
	})
	if err != nil {
		return ClassifiedVariety{}, fmt.Errorf("openai classify variety: %w", err)
	}

	if len(resp.Choices) == 0 {
		return ClassifiedVariety{}, fmt.Errorf("openai returned no choices")
	}

	var result ClassifiedVariety
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return ClassifiedVariety{}, fmt.Errorf("parse classify response: %w", err)
	}

	return result, nil
}

// BackfillUnclassified fetches coffees where variety_raw is set but variety is
// not, runs the LLM classifier, and updates the DB. Returns counts of
// classified and failed rows.
func (c *Classifier) BackfillUnclassified(ctx context.Context, queries *db.Queries) (classified, failed int) {
	rows, err := queries.ListCoffeesNeedingVariety(ctx)
	if err != nil {
		slog.Error("list coffees needing variety", "error", err)
		return 0, 0
	}

	if len(rows) == 0 {
		return 0, 0
	}

	slog.Info("coffees needing variety classification", "count", len(rows))

	for _, row := range rows {
		result, err := c.ClassifyVariety(ctx, row.VarietyRaw.String)
		if err != nil {
			slog.Warn("classify variety failed", "id", row.ID, "variety_raw", row.VarietyRaw.String, "error", err)
			failed++
			continue
		}

		if result.Variety == "" {
			slog.Debug("variety unrecognised by LLM", "id", row.ID, "variety_raw", row.VarietyRaw.String)
			failed++
			continue
		}

		err = queries.UpdateCoffeeVariety(ctx, db.UpdateCoffeeVarietyParams{
			ID:      row.ID,
			Variety: pgtype.Text{String: result.Variety, Valid: true},
			Species: pgtype.Text{String: result.Species, Valid: result.Species != ""},
		})
		if err != nil {
			slog.Warn("update coffee variety", "id", row.ID, "error", err)
			failed++
			continue
		}

		slog.Info("classified variety", "id", row.ID, "raw", row.VarietyRaw.String,
			"variety", result.Variety, "species", result.Species)
		classified++

		time.Sleep(500 * time.Millisecond)
	}

	return classified, failed
}
