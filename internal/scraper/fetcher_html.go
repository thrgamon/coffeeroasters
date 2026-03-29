package scraper

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// FetchHTML fetches a roaster's shop page, cleans the HTML, and uses the LLM
// to extract coffee product listings.
func FetchHTMLPage(ctx context.Context, cfg domain.RoasterConfig, client *http.Client, extractor *Extractor, limiter *RateLimiter) ([]RawCoffee, error) {
	limiter.Wait(cfg.ShopURL, 3*time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ShopURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", botUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", cfg.ShopURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", cfg.ShopURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	cleaned, err := CleanHTML(string(body), cfg.ContentSelector)
	if err != nil {
		return nil, fmt.Errorf("clean HTML: %w", err)
	}

	slog.Info("extracting from HTML page", "roaster", cfg.Slug, "html_size", len(cleaned))

	products, err := extractor.ExtractFromPage(ctx, cleaned)
	if err != nil {
		return nil, fmt.Errorf("page extraction: %w", err)
	}

	var skipped int
	for _, p := range products {
		if !p.IsCoffee {
			skipped++
			slog.Debug("skipped non-coffee product", "roaster", cfg.Slug, "name", p.Name)
		}
	}
	if skipped > 0 {
		slog.Info("filtered non-coffee products", "roaster", cfg.Slug, "total", len(products), "skipped", skipped)
	}

	baseURL, _ := url.Parse(cfg.ShopURL)

	var coffees []RawCoffee
	for _, p := range products {
		if !p.IsCoffee {
			continue
		}
		raw := RawCoffee{
			Name:      p.Name,
			InStock:   p.InStock,
			ScrapedAt: time.Now(),
			Currency:  "AUD",
		}

		if p.ProductURL != nil {
			raw.ProductURL = resolveURL(baseURL, *p.ProductURL)
		}
		if p.PriceText != nil {
			raw.PriceRaw = *p.PriceText
		}
		if p.WeightText != nil {
			raw.WeightRaw = *p.WeightText
		}
		if p.Origin != nil {
			raw.OriginRaw = *p.Origin
		}
		if p.Region != nil {
			raw.RegionRaw = *p.Region
		}
		if p.Process != nil {
			raw.ProcessRaw = *p.Process
		}
		if p.RoastLevel != nil {
			raw.RoastRaw = *p.RoastLevel
		}
		if p.TastingNotes != nil {
			raw.TastingNotes = *p.TastingNotes
		}
		if p.Variety != nil {
			raw.VarietyRaw = *p.Variety
		}
		if p.Producer != nil {
			raw.ProducerRaw = *p.Producer
		}
		if p.Description != nil {
			raw.Description = *p.Description
		}

		raw.IsBlend = p.IsBlend
		if p.IsBlend && len(p.BlendComponents) > 0 {
			for _, bc := range p.BlendComponents {
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

		coffees = append(coffees, raw)
	}

	return coffees, nil
}

func resolveURL(base *url.URL, href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}
