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

// BlendComponent represents one origin component of a blend coffee, as
// extracted by the LLM.
type BlendComponent struct {
	Origin     *string `json:"origin" jsonschema:"description=Country of origin for this blend component. Null if not stated."`
	Region     *string `json:"region" jsonschema:"description=Growing region for this blend component. Null if not stated."`
	Variety    *string `json:"variety" jsonschema:"description=Coffee variety for this blend component. Null if not stated."`
	Percentage *int    `json:"percentage" jsonschema:"description=Percentage of this component in the blend. Null if not stated."`
}

// BatchProduct pairs an index with the extracted coffee data. The index
// corresponds to the position in the input descriptions array.
type BatchProduct struct {
	Index           int              `json:"index" jsonschema:"description=Zero-based index matching the input description array"`
	IsCoffee        bool             `json:"is_coffee" jsonschema:"description=True only if this is a single purchasable bag of whole coffee beans. False for ground coffee, drip bags, instant coffee, cold brew, ready-to-drink, concentrates, capsules, pods, subscriptions, recurring plans, sample packs, bundles, gift cards, equipment, merchandise, apparel, tea, chocolate, vouchers, courses, and accessories."`
	Name            string           `json:"name" jsonschema:"description=Product name"`
	Origin          *string          `json:"origin" jsonschema:"description=Country or region of origin. Null if not stated. For blends with multiple origins use the blend_components array instead."`
	Region          *string          `json:"region" jsonschema:"description=Specific growing region. Null if not stated."`
	Process         *string          `json:"process" jsonschema:"description=Processing method. Null if not stated."`
	RoastLevel      *string          `json:"roast_level" jsonschema:"description=Roast level. Null if not stated."`
	TastingNotes    *string          `json:"tasting_notes" jsonschema:"description=Comma-separated tasting notes. Null if not stated."`
	Variety         *string          `json:"variety" jsonschema:"description=Coffee variety. Null if not stated."`
	Producer        *string          `json:"producer" jsonschema:"description=Farm, estate, cooperative, or washing station name. Null if not stated."`
	Description     *string          `json:"description" jsonschema:"description=The roaster's flavour writeup or descriptive text about this coffee. Copy the relevant prose as-is. Null if no descriptive text beyond structured data."`
	BrewRecipe      *string          `json:"brew_recipe" jsonschema:"description=Any brewing recommendations, brew guide, recipe, or preparation tips provided by the roaster. Include method, dose, ratio, temperature, grind size, brew time, or any other brewing instructions. Copy the text as-is. Null if no brewing information is stated."`
	IsBlend         bool             `json:"is_blend" jsonschema:"description=True if this coffee is a blend of beans from multiple origins or farms."`
	BlendComponents []BlendComponent `json:"blend_components" jsonschema:"description=For blends only: the individual origin components. Empty array if not a blend or components are not stated."`
}

// PageExtraction is the structured output for extracting coffees from a
// full HTML page (non-Shopify path).
type PageExtraction struct {
	Products []PageProduct `json:"products"`
}

// PageProduct extends CoffeeProduct with pricing and URL data that must
// be extracted from the HTML (since there is no Shopify JSON to provide it).
type PageProduct struct {
	IsCoffee        bool             `json:"is_coffee" jsonschema:"description=True only if this is a single purchasable bag of whole coffee beans. False for ground coffee, drip bags, instant coffee, cold brew, ready-to-drink, concentrates, capsules, pods, subscriptions, recurring plans, sample packs, bundles, gift cards, equipment, merchandise, apparel, tea, chocolate, vouchers, courses, and accessories."`
	Name            string           `json:"name" jsonschema:"description=Product name"`
	Origin          *string          `json:"origin" jsonschema:"description=Country or region of origin. Null if not stated. For blends with multiple origins use the blend_components array instead."`
	Region          *string          `json:"region" jsonschema:"description=Specific growing region. Null if not stated."`
	Process         *string          `json:"process" jsonschema:"description=Processing method. Null if not stated."`
	RoastLevel      *string          `json:"roast_level" jsonschema:"description=Roast level. Null if not stated."`
	TastingNotes    *string          `json:"tasting_notes" jsonschema:"description=Comma-separated tasting notes. Null if not stated."`
	Variety         *string          `json:"variety" jsonschema:"description=Coffee variety. Null if not stated."`
	Producer        *string          `json:"producer" jsonschema:"description=Farm, estate, cooperative, or washing station name. Null if not stated."`
	Description     *string          `json:"description" jsonschema:"description=The roaster's flavour writeup or descriptive text about this coffee. Copy the relevant prose as-is. Null if no descriptive text beyond structured data."`
	BrewRecipe      *string          `json:"brew_recipe" jsonschema:"description=Any brewing recommendations, brew guide, recipe, or preparation tips provided by the roaster. Include method, dose, ratio, temperature, grind size, brew time, or any other brewing instructions. Copy the text as-is. Null if no brewing information is stated."`
	PriceText       *string          `json:"price_text" jsonschema:"description=Price as displayed (e.g. $32.00). Null if not found."`
	WeightText      *string          `json:"weight_text" jsonschema:"description=Weight as displayed (e.g. 250g). Null if not found."`
	InStock         bool             `json:"in_stock" jsonschema:"description=Whether the product appears to be in stock"`
	ProductURL      *string          `json:"product_url" jsonschema:"description=URL or relative path to the product page. Null if not found."`
	IsBlend         bool             `json:"is_blend" jsonschema:"description=True if this coffee is a blend of beans from multiple origins or farms."`
	BlendComponents []BlendComponent `json:"blend_components" jsonschema:"description=For blends only: the individual origin components. Empty array if not a blend or components are not stated."`
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

Only extract products that are a single purchasable bag of whole coffee beans.
Set is_coffee=true for these. Set is_coffee=false for everything else including: ground coffee, pre-ground coffee, drip bags, instant coffee, cold brew, ready-to-drink beverages, coffee concentrates, capsules, pods, subscriptions, recurring plans, sample packs, mixed bundles, gift cards, equipment, grinders, merchandise, apparel, tea, chocolate, vouchers, accessories, cleaning products, and courses.

For each coffee product, extract:
- name: the product name
- origin: country or region of origin (for single-origin coffees)
- region: specific growing region within the country
- process: processing method (Washed, Natural, Honey, Anaerobic, etc.)
- roast_level: roast profile (Light, Medium, Filter, Espresso, etc.)
- tasting_notes: comma-separated flavour descriptors
- variety: coffee cultivar/variety
- producer: farm, estate, cooperative, or washing station name
- description: the roaster's descriptive prose about this coffee (flavour writeups, producer story). Copy the text as-is. Null if none present.
- brew_recipe: any brewing recommendations, brew guide, recipe, or preparation tips (e.g. dose, water ratio, temperature, grind size, brew time, method). Copy the text as-is. Null if no brewing information is stated.
- is_blend: true if the coffee is a blend of beans from multiple origins or farms
- blend_components: for blends, list each origin component with country, region, variety, and percentage if stated

For blends (e.g. "Seasonal Espresso Blend" with origins "Colombia, Brazil, Guatemala"), set is_blend=true and populate blend_components with each origin. Leave the top-level origin field null for blends.

Use null for any field where the information is not clearly stated. Do not guess.
Prices are in AUD unless stated otherwise.`

const detailPageSystemPrompt = `You extract detailed coffee product information from an individual product page on an Australian specialty coffee roaster website.

This is a single product page, not a listing of many products. Extract as much detail as possible.

For the coffee product on this page, extract:
- name: the product name
- origin: country of origin (e.g. Colombia, Ethiopia, Kenya). For blends with multiple origins, leave null and use blend_components instead.
- region: specific growing region within the country (e.g. Huila, Yirgacheffe, Nyeri)
- process: processing method (e.g. Washed, Natural, Honey, Anaerobic)
- roast_level: roast profile (e.g. Light, Medium, Filter, Espresso, Omni)
- tasting_notes: comma-separated flavour descriptors
- variety: coffee cultivar/variety (e.g. Caturra, SL28, Gesha, Heirloom, Castillo)
- producer: farm name, estate name, cooperative, or washing station name
- description: the roaster's descriptive prose about this coffee (flavour writeups, producer story). Copy the text as-is. Null if none present.
- brew_recipe: any brewing recommendations, brew guide, recipe, or preparation tips (e.g. dose, water ratio, temperature, grind size, brew time, method). Copy the text as-is. Null if no brewing information is stated.
- price_text: the displayed price (e.g. "$32.00")
- weight_text: the displayed weight (e.g. "250g")
- in_stock: whether it appears available for purchase
- is_blend: true if the coffee is a blend of beans from multiple origins or farms
- blend_components: for blends, list each origin component with country, region, variety, and percentage if stated

Use null for any field where the information is not clearly stated. Do not guess.
Only extract if this is a single purchasable bag of whole coffee beans. Set is_coffee=true for these.
Set is_coffee=false for everything else: ground coffee, pre-ground coffee, drip bags, instant coffee, cold brew, ready-to-drink beverages, coffee concentrates, capsules, pods, subscriptions, recurring plans, sample packs, mixed bundles, gift cards, equipment, grinders, merchandise, apparel, tea, chocolate, vouchers, accessories, cleaning products, and courses.`

const pageSystemPrompt = `You extract coffee product listings from Australian specialty coffee roaster web pages.

For each coffee product you find on the page, extract:
- name: the product name
- origin, region, process, roast_level, tasting_notes, variety, producer: as described
- description: the roaster's descriptive prose about this coffee (flavour writeups, producer story). Copy the text as-is. Null if none present.
- brew_recipe: any brewing recommendations, brew guide, recipe, or preparation tips (e.g. dose, water ratio, temperature, grind size, brew time, method). Copy the text as-is. Null if no brewing information is stated.
- price_text: the displayed price (e.g. "$32.00")
- weight_text: the displayed weight (e.g. "250g")
- in_stock: whether it appears available for purchase
- product_url: the link to the product page (relative or absolute URL)
- is_blend: true if the coffee is a blend of beans from multiple origins or farms
- blend_components: for blends, list each origin component with country, region, variety, and percentage if stated

For blends, leave the top-level origin field null and populate blend_components instead.

Use null for any field where the information is not clearly stated. Do not guess.
Only extract products that are a single purchasable bag of whole coffee beans. Set is_coffee=true for these.
Set is_coffee=false for everything else: ground coffee, pre-ground coffee, drip bags, instant coffee, cold brew, ready-to-drink beverages, coffee concentrates, capsules, pods, subscriptions, recurring plans, sample packs, mixed bundles, gift cards, equipment, grinders, merchandise, apparel, tea, chocolate, vouchers, accessories, cleaning products, and courses.`

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

// ExtractFromProductPage sends a single product page's cleaned HTML to
// GPT-4.1-mini and returns the extracted coffee data. Returns nil if the
// page does not contain a coffee product.
func (e *Extractor) ExtractFromProductPage(ctx context.Context, cleanedHTML string) (*PageProduct, error) {
	schema := llm.GenerateSchema(PageExtraction{})

	resp, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(detailPageSystemPrompt),
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
		return nil, fmt.Errorf("openai detail page extraction: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	var result PageExtraction
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("parse detail page extraction response: %w", err)
	}

	if len(result.Products) == 0 {
		return nil, nil
	}

	return &result.Products[0], nil
}

