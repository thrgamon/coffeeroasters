package scraper

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"

	"github.com/thrgamon/coffeeroasters/internal/llm"
)

// CoffeeProduct is the structured output schema for a single coffee product
// extracted by the LLM. Pointer types indicate optional/nullable fields.
type CoffeeProduct struct {
	Name         string   `json:"name" jsonschema:"description=Product name as shown on the page"`
	Origin       *string  `json:"origin" jsonschema:"description=Country or region of origin (e.g. Ethiopia or Colombia Huila). Null if not stated."`
	Region       *string  `json:"region" jsonschema:"description=Specific growing region (e.g. Yirgacheffe or Gedeo Zone). Null if not stated."`
	Process      *string  `json:"process" jsonschema:"description=Processing method (e.g. Washed or Natural or Honey). Null if not stated."`
	RoastLevel   *string  `json:"roast_level" jsonschema:"description=Roast level (e.g. Light or Medium or Filter or Espresso). Null if not stated."`
	TastingNotes *string  `json:"tasting_notes" jsonschema:"description=Comma-separated tasting/flavour notes. Null if not stated."`
	Variety      *string  `json:"variety" jsonschema:"description=Coffee variety (e.g. Heirloom or SL28 or Caturra). Null if not stated."`
}

// BatchExtraction is the structured output for extracting multiple coffees
// from Shopify product descriptions (batched into a single LLM call).
type BatchExtraction struct {
	Products []BatchProduct `json:"products"`
}

// BatchProduct pairs an index with the extracted coffee data. The index
// corresponds to the position in the input descriptions array.
type BatchProduct struct {
	Index        int     `json:"index" jsonschema:"description=Zero-based index matching the input description array"`
	IsCoffee     bool    `json:"is_coffee" jsonschema:"description=True only if this is a single purchasable bag of coffee beans or ground coffee. False for subscriptions and recurring plans and sample packs and bundles and gift cards and equipment and merchandise and tea and vouchers and courses and accessories."`
	Name         string  `json:"name" jsonschema:"description=Product name"`
	Origin       *string `json:"origin" jsonschema:"description=Country or region of origin. Null if not stated."`
	Region       *string `json:"region" jsonschema:"description=Specific growing region. Null if not stated."`
	Process      *string `json:"process" jsonschema:"description=Processing method. Null if not stated."`
	RoastLevel   *string `json:"roast_level" jsonschema:"description=Roast level. Null if not stated."`
	TastingNotes *string `json:"tasting_notes" jsonschema:"description=Comma-separated tasting notes. Null if not stated."`
	Variety      *string `json:"variety" jsonschema:"description=Coffee variety. Null if not stated."`
	Producer     *string `json:"producer" jsonschema:"description=Farm, estate, cooperative, or washing station name. Null if not stated."`
}

// PageExtraction is the structured output for extracting coffees from a
// full HTML page (non-Shopify path).
type PageExtraction struct {
	Products []PageProduct `json:"products"`
}

// PageProduct extends CoffeeProduct with pricing and URL data that must
// be extracted from the HTML (since there is no Shopify JSON to provide it).
type PageProduct struct {
	IsCoffee     bool    `json:"is_coffee" jsonschema:"description=True only if this is a single purchasable bag of coffee beans or ground coffee. False for subscriptions and recurring plans and sample packs and bundles and gift cards and equipment and merchandise and tea and vouchers and courses and accessories."`
	Name         string  `json:"name" jsonschema:"description=Product name"`
	Origin       *string `json:"origin" jsonschema:"description=Country or region of origin. Null if not stated."`
	Region       *string `json:"region" jsonschema:"description=Specific growing region. Null if not stated."`
	Process      *string `json:"process" jsonschema:"description=Processing method. Null if not stated."`
	RoastLevel   *string `json:"roast_level" jsonschema:"description=Roast level. Null if not stated."`
	TastingNotes *string `json:"tasting_notes" jsonschema:"description=Comma-separated tasting notes. Null if not stated."`
	Variety      *string `json:"variety" jsonschema:"description=Coffee variety. Null if not stated."`
	Producer     *string `json:"producer" jsonschema:"description=Farm, estate, cooperative, or washing station name. Null if not stated."`
	PriceText    *string `json:"price_text" jsonschema:"description=Price as displayed (e.g. $32.00). Null if not found."`
	WeightText   *string `json:"weight_text" jsonschema:"description=Weight as displayed (e.g. 250g). Null if not found."`
	InStock      bool    `json:"in_stock" jsonschema:"description=Whether the product appears to be in stock"`
	ProductURL   *string `json:"product_url" jsonschema:"description=URL or relative path to the product page. Null if not found."`
}

// ProductDescription is an input to the batch extraction: a Shopify product
// title paired with its HTML description body.
type ProductDescription struct {
	Index int
	Title string
	HTML  string
}

// Extractor uses GPT-4.1-mini with structured outputs to extract coffee
// product data from roaster websites.
type Extractor struct {
	client *openai.Client
}

// NewExtractor creates an Extractor. The OpenAI client should be configured
// with the OPENAI_API_KEY environment variable.
func NewExtractor(client *openai.Client) *Extractor {
	return &Extractor{client: client}
}

const batchSystemPrompt = `You extract coffee product information from Australian specialty coffee roaster websites.

Only extract products that are a single purchasable bag of coffee beans or ground coffee.
Set is_coffee=true for these. Set is_coffee=false for everything else including: subscriptions, recurring plans, sample packs, mixed bundles, gift cards, equipment, merchandise, apparel, tea, chocolate, vouchers, accessories, cleaning products, courses, drip bags, and pods.

For each coffee product, extract:
- name: the product name
- origin: country or region of origin
- region: specific growing region within the country
- process: processing method (Washed, Natural, Honey, Anaerobic, etc.)
- roast_level: roast profile (Light, Medium, Filter, Espresso, etc.)
- tasting_notes: comma-separated flavour descriptors
- variety: coffee cultivar/variety
- producer: farm, estate, cooperative, or washing station name

Use null for any field where the information is not clearly stated. Do not guess.
Prices are in AUD unless stated otherwise.`

const pageSystemPrompt = `You extract coffee product listings from Australian specialty coffee roaster web pages.

For each coffee product you find on the page, extract:
- name: the product name
- origin, region, process, roast_level, tasting_notes, variety, producer: as described
- price_text: the displayed price (e.g. "$32.00")
- weight_text: the displayed weight (e.g. "250g")
- in_stock: whether it appears available for purchase
- product_url: the link to the product page (relative or absolute URL)

Use null for any field where the information is not clearly stated. Do not guess.
Only extract products that are a single purchasable bag of coffee beans or ground coffee. Set is_coffee=true for these.
Set is_coffee=false for everything else: subscriptions, recurring plans, sample packs, mixed bundles, gift cards, equipment, merchandise, tea, vouchers, accessories, pods, and drip bags.`

// ExtractFromDescriptions sends a batch of Shopify product descriptions to
// GPT-4.1-mini and returns structured coffee data for each.
func (e *Extractor) ExtractFromDescriptions(ctx context.Context, descriptions []ProductDescription) ([]BatchProduct, error) {
	if len(descriptions) == 0 {
		return nil, nil
	}

	userContent := "Extract coffee information from these product descriptions:\n\n"
	for _, d := range descriptions {
		userContent += fmt.Sprintf("--- Product %d: %s ---\n%s\n\n", d.Index, d.Title, d.HTML)
	}

	schema := llm.GenerateSchema(BatchExtraction{})

	resp, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(batchSystemPrompt),
			openai.UserMessage(userContent),
		},
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "batch_extraction",
					Strict: openai.Opt(true),
					Schema: schema,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai batch extraction: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	var result BatchExtraction
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("parse batch extraction response: %w", err)
	}

	return result.Products, nil
}

// ExtractFromPage sends cleaned HTML to GPT-4.1-mini and returns structured
// coffee product data including pricing and URLs.
func (e *Extractor) ExtractFromPage(ctx context.Context, cleanedHTML string) ([]PageProduct, error) {
	schema := llm.GenerateSchema(PageExtraction{})

	resp, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(pageSystemPrompt),
			openai.UserMessage(cleanedHTML),
		},
		Temperature: openai.Opt(0.0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "page_extraction",
					Strict: openai.Opt(true),
					Schema: schema,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai page extraction: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	raw := resp.Choices[0].Message.Content

	var result PageExtraction
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse page extraction response: %w", err)
	}

	return result.Products, nil
}

