package embedding

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	openai "github.com/openai/openai-go/v2"

	"github.com/thrgamon/coffeeroasters/internal/db"
)

// Embedder generates text embeddings via OpenAI's text-embedding-3-small model.
type Embedder struct {
	client *openai.Client
}

// NewEmbedder creates an Embedder from an OpenAI client.
func NewEmbedder(client *openai.Client) *Embedder {
	return &Embedder{client: client}
}

// Embed returns the embedding vector for the given text.
func (e *Embedder) Embed(ctx context.Context, text string) ([]float64, error) {
	resp, err := e.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Small,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("openai returned no embeddings")
	}

	return resp.Data[0].Embedding, nil
}

// BackfillPending embeds all coffees that have a description but no embedding.
// Returns counts of embedded and failed coffees.
func (e *Embedder) BackfillPending(ctx context.Context, queries *db.Queries) (embedded, failed int) {
	rows, err := queries.ListCoffeesNeedingEmbedding(ctx)
	if err != nil {
		slog.Error("list coffees needing embedding", "error", err)
		return 0, 0
	}

	if len(rows) == 0 {
		return 0, 0
	}

	slog.Info("coffees needing embedding", "count", len(rows))

	for _, row := range rows {
		vec, err := e.Embed(ctx, row.Description.String)
		if err != nil {
			slog.Warn("embedding failed", "coffee_id", row.ID, "error", err)
			failed++
			continue
		}

		err = queries.UpdateCoffeeEmbedding(ctx, db.UpdateCoffeeEmbeddingParams{
			ID:        row.ID,
			Embedding: vec,
		})
		if err != nil {
			slog.Warn("update embedding", "coffee_id", row.ID, "error", err)
			failed++
			continue
		}

		slog.Info("embedded coffee", "coffee_id", row.ID)
		embedded++

		time.Sleep(100 * time.Millisecond)
	}

	return embedded, failed
}

// CosineSimilarity computes the cosine similarity between two vectors.
// Returns 0 for nil, empty, or mismatched-length vectors.
func CosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}

	return dot / denom
}
