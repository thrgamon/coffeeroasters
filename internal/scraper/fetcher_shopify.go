package scraper

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thrgamon/coffeeroasters/internal/domain"
	"github.com/thrgamon/coffeeroasters/internal/normalise"
)

// shopifyProduct mirrors the relevant fields from Shopify's /products.json.
type shopifyProduct struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	BodyHTML  string           `json:"body_html"`
	Handle    string           `json:"handle"`
	Type      string           `json:"product_type"`
	Variants  []shopifyVariant `json:"variants"`
	Images    []shopifyImage   `json:"images"`
	Status    string           `json:"status"`
	CreatedAt string           `json:"created_at"`
}

type shopifyVariant struct {
	Price      string  `json:"price"`
	Available  bool    `json:"available"`
	Title      string  `json:"title"`
	Grams      int     `json:"grams"`
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weight_unit"`
}

type shopifyImage struct {
	Src string `json:"src"`
}

type shopifyResponse struct {
	Products []shopifyProduct `json:"products"`
}

const shopifyBatchSize = 10

// ShopifyFetchResult separates changed products (needing LLM extraction) from
// unchanged ones (only price/stock updates needed).
type ShopifyFetchResult struct {
	Changed   []RawCoffee // Products with new/changed content, fully extracted
	Unchanged []RawCoffee // Products with unchanged content, only Shopify JSON fields populated
}

// FetchShopify fetches products from a Shopify store's /products.json endpoint,
// extracts structured data from the JSON, and uses the LLM to parse origin/
// process/tasting notes from the HTML descriptions.
//
// knownHashes maps product_url -> source_hash for previously scraped products.
// Products whose body_html hash matches a known hash skip LLM extraction.
func FetchShopify(ctx context.Context, cfg domain.RoasterConfig, client *http.Client, extractor *Extractor, limiter *RateLimiter, knownHashes map[string]string) (*ShopifyFetchResult, error) {
	var allProducts []shopifyProduct

	baseURL := strings.TrimSuffix(cfg.ShopURL, "/")
	// Ensure we're hitting the root /products.json, not a collection
	productsURL := baseURL
	if !strings.HasSuffix(productsURL, "/products.json") {
		// Strip any collection path and use root products.json
		parts := strings.Split(productsURL, "/")
		scheme := parts[0] + "//" + parts[2]
		productsURL = scheme + "/products.json"
	}

	for page := 1; page <= 5; page++ {
		url := fmt.Sprintf("%s?limit=250&page=%d", productsURL, page)
		limiter.Wait(url, 3*time.Second)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}
		req.Header.Set("User-Agent", botUserAgent)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GET %s: %w", url, err)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
		}
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		var sr shopifyResponse
		if err := json.Unmarshal(body, &sr); err != nil {
			return nil, fmt.Errorf("parse products.json: %w", err)
		}

		if len(sr.Products) == 0 {
			break
		}
		allProducts = append(allProducts, sr.Products...)

		if len(sr.Products) < 250 {
			break
		}
	}

	// Filter by product_type if configured, and skip archived/upcoming products
	var filtered []shopifyProduct
	for _, p := range allProducts {
		if cfg.ProductType != "" && !strings.EqualFold(p.Type, cfg.ProductType) {
			continue
		}
		typeLower := strings.ToLower(p.Type)
		if strings.Contains(typeLower, "archive") || strings.Contains(typeLower, "upcoming") {
			continue
		}
		filtered = append(filtered, p)
	}

	if len(filtered) == 0 {
		return &ShopifyFetchResult{}, nil
	}

	// Separate changed vs unchanged products based on body_html hash
	var changed []shopifyProduct
	result := &ShopifyFetchResult{}

	for _, p := range filtered {
		hash := hashBody(p.BodyHTML)
		productURL := shopifyProductURL(cfg.Website, p.Handle)

		if existing, ok := knownHashes[productURL]; ok && existing == hash {
			// Unchanged: populate only Shopify JSON fields (no LLM extraction)
			raw := shopifyJSONFields(cfg, p, hash)
			result.Unchanged = append(result.Unchanged, raw)
		} else {
			changed = append(changed, p)
		}
	}

	if len(changed) > 0 {
		slog.Info("change detection",
			"roaster", cfg.Slug,
			"total", len(filtered),
			"changed", len(changed),
			"unchanged", len(result.Unchanged))
	}

	// Extract changed products via LLM in batches
	for i := 0; i < len(changed); i += shopifyBatchSize {
		end := i + shopifyBatchSize
		if end > len(changed) {
			end = len(changed)
		}
		batch := changed[i:end]

		var descriptions []ProductDescription
		for j, p := range batch {
			descriptions = append(descriptions, ProductDescription{
				Index: j,
				Title: p.Title,
				HTML:  p.BodyHTML,
			})
		}

		slog.Info("extracting batch", "roaster", cfg.Slug, "batch_start", i, "batch_size", len(batch))

		extracted, err := extractor.ExtractFromDescriptions(ctx, descriptions)
		if err != nil {
			slog.Error("batch extraction failed", "roaster", cfg.Slug, "error", err)
			continue
		}

		// Merge LLM extraction with Shopify JSON data
		for _, ep := range extracted {
			if ep.Index < 0 || ep.Index >= len(batch) {
				continue
			}
			if !ep.IsCoffee {
				continue
			}
			sp := batch[ep.Index]

			raw := shopifyJSONFields(cfg, sp, hashBody(sp.BodyHTML))
			raw.Name = ep.Name

			// Origin, process, tasting notes from LLM
			if ep.Origin != nil {
				raw.OriginRaw = *ep.Origin
			}
			if ep.Region != nil {
				raw.RegionRaw = *ep.Region
			}
			if ep.Process != nil {
				raw.ProcessRaw = *ep.Process
			}
			if ep.RoastLevel != nil {
				raw.RoastRaw = *ep.RoastLevel
			}
			if ep.TastingNotes != nil {
				raw.TastingNotes = *ep.TastingNotes
			}
			if ep.Variety != nil {
				raw.VarietyRaw = *ep.Variety
			}
			if ep.Producer != nil {
				raw.ProducerRaw = *ep.Producer
			}
			if ep.Description != nil {
				raw.Description = *ep.Description
			}

			raw.IsBlend = ep.IsBlend
			if ep.IsBlend && len(ep.BlendComponents) > 0 {
				for _, bc := range ep.BlendComponents {
					comp := RawBlendComponent{}
					if bc.Origin != nil {
						comp.Origin = *bc.Origin
					}
					if bc.Region != nil {
						comp.Region = *bc.Region
					}
					if bc.Variety != nil {
						comp.Variety = *bc.Variety
					}
					if bc.Percentage != nil {
						comp.Percentage = *bc.Percentage
					}
					raw.BlendComponents = append(raw.BlendComponents, comp)
				}
			}

			result.Changed = append(result.Changed, raw)
		}
	}

	return result, nil
}

// shopifyJSONFields builds a RawCoffee populated only from Shopify JSON data
// (price, stock, weight, image). No LLM extraction fields.
func shopifyJSONFields(cfg domain.RoasterConfig, sp shopifyProduct, hash string) RawCoffee {
	raw := RawCoffee{
		Name:       sp.Title,
		ProductURL: shopifyProductURL(cfg.Website, sp.Handle),
		SourceHash: hash,
		ScrapedAt:  time.Now(),
		Currency:   "AUD",
	}

	if len(sp.Images) > 0 {
		raw.ImageURL = sp.Images[0].Src
	}

	if len(sp.Variants) > 0 {
		var bestPer100g int64 = math.MaxInt64
		var anyInStock bool
		var minPer100g, maxPer100g int64

		for _, v := range sp.Variants {
			if v.Available {
				anyInStock = true
			}

			priceCents, priceOK := normalise.NormalisePriceAUD(v.Price)
			weightGrams, weightOK := normalise.NormaliseWeightGrams(shopifyWeight(v))
			if !priceOK || !weightOK || weightGrams <= 0 {
				continue
			}

			per100g := normalise.PricePer100g(priceCents, weightGrams)
			if per100g <= 0 {
				continue
			}

			if minPer100g == 0 || per100g < minPer100g {
				minPer100g = per100g
			}
			if per100g > maxPer100g {
				maxPer100g = per100g
			}

			if per100g < bestPer100g {
				bestPer100g = per100g
				raw.PriceRaw = v.Price
				raw.WeightRaw = shopifyWeight(v)
			}
		}

		raw.InStock = anyInStock
		raw.PricePer100gMin = minPer100g
		raw.PricePer100gMax = maxPer100g

		if bestPer100g == math.MaxInt64 {
			raw.PriceRaw = sp.Variants[0].Price
			raw.WeightRaw = shopifyWeight(sp.Variants[0])
			raw.InStock = sp.Variants[0].Available
		}
	}

	return raw
}

func hashBody(body string) string {
	h := sha256.Sum256([]byte(body))
	return fmt.Sprintf("%x", h)
}

func shopifyProductURL(website, handle string) string {
	base := strings.TrimSuffix(website, "/")
	return base + "/products/" + handle
}

func shopifyWeight(v shopifyVariant) string {
	if v.Grams > 0 {
		return strconv.Itoa(v.Grams) + "g"
	}
	if v.Weight > 0 {
		unit := v.WeightUnit
		if unit == "" {
			unit = "g"
		}
		return fmt.Sprintf("%.0f%s", v.Weight, unit)
	}
	return ""
}
