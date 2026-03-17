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

	"github.com/PuerkitoBio/goquery"

	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// FetchHTMLDetail implements a two-pass fetch: first it scrapes the listing
// page for product URLs, then fetches each product page individually and
// extracts detailed coffee data via the LLM.
func FetchHTMLDetail(ctx context.Context, cfg domain.RoasterConfig, client *http.Client, extractor *Extractor, limiter *RateLimiter) ([]RawCoffee, error) {
	// Pass 1: fetch listing page and extract product URLs.
	productURLs, err := extractProductURLs(ctx, cfg, client, limiter)
	if err != nil {
		return nil, fmt.Errorf("extract product URLs: %w", err)
	}

	slog.Info("found product URLs", "roaster", cfg.Slug, "count", len(productURLs))

	if len(productURLs) == 0 {
		return nil, nil
	}

	// Pass 2: fetch each product page and extract coffee details.
	var coffees []RawCoffee
	for i, productURL := range productURLs {
		slog.Info("extracting product", "roaster", cfg.Slug, "progress", fmt.Sprintf("%d/%d", i+1, len(productURLs)), "url", productURL)

		raw, err := fetchAndExtractProduct(ctx, cfg, client, extractor, limiter, productURL)
		if err != nil {
			slog.Warn("product extraction failed", "roaster", cfg.Slug, "url", productURL, "error", err)
			continue
		}
		if raw == nil {
			slog.Debug("skipped non-coffee product", "roaster", cfg.Slug, "url", productURL)
			continue
		}

		coffees = append(coffees, *raw)
	}

	return coffees, nil
}

// extractProductURLs fetches the listing page and extracts deduplicated,
// same-domain product links.
func extractProductURLs(ctx context.Context, cfg domain.RoasterConfig, client *http.Client, limiter *RateLimiter) ([]string, error) {
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", cfg.ShopURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	// Scope to content selector if configured.
	var selection *goquery.Selection
	if cfg.ContentSelector != "" {
		selection = doc.Find(cfg.ContentSelector)
	} else {
		selection = doc.Find("body")
	}

	baseURL, err := url.Parse(cfg.ShopURL)
	if err != nil {
		return nil, fmt.Errorf("parse shop URL: %w", err)
	}

	seen := make(map[string]bool)
	var urls []string

	selection.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		resolved := resolveURL(baseURL, href)

		// Filter to same-domain links only.
		parsed, err := url.Parse(resolved)
		if err != nil {
			return
		}
		if parsed.Host != baseURL.Host {
			return
		}

		// Skip the listing page itself, anchors, and non-path links.
		if resolved == cfg.ShopURL || parsed.Path == baseURL.Path {
			return
		}

		if seen[resolved] {
			return
		}
		seen[resolved] = true
		urls = append(urls, resolved)
	})

	return urls, nil
}

// fetchAndExtractProduct fetches a single product page, cleans the HTML,
// and extracts coffee data via the LLM. Returns nil if not a coffee product.
func fetchAndExtractProduct(ctx context.Context, cfg domain.RoasterConfig, client *http.Client, extractor *Extractor, limiter *RateLimiter, productURL string) (*RawCoffee, error) {
	limiter.Wait(productURL, 3*time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, productURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", botUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", productURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", productURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// Use detail_selector if configured, otherwise no selector.
	cleaned, err := CleanHTML(string(body), cfg.DetailSelector)
	if err != nil {
		return nil, fmt.Errorf("clean HTML: %w", err)
	}

	slog.Debug("cleaned product page", "url", productURL, "html_size", len(cleaned))

	product, err := extractor.ExtractFromProductPage(ctx, cleaned)
	if err != nil {
		return nil, fmt.Errorf("extract product: %w", err)
	}

	if product == nil || !product.IsCoffee {
		return nil, nil
	}

	raw := RawCoffee{
		Name:       product.Name,
		ProductURL: productURL,
		InStock:    product.InStock,
		ScrapedAt:  time.Now(),
		Currency:   "AUD",
	}

	if product.PriceText != nil {
		raw.PriceRaw = *product.PriceText
	}
	if product.WeightText != nil {
		raw.WeightRaw = *product.WeightText
	}
	if product.Origin != nil {
		raw.OriginRaw = *product.Origin
	}
	if product.Region != nil {
		raw.RegionRaw = *product.Region
	}
	if product.Process != nil {
		raw.ProcessRaw = *product.Process
	}
	if product.RoastLevel != nil {
		raw.RoastRaw = *product.RoastLevel
	}
	if product.TastingNotes != nil {
		raw.TastingNotes = *product.TastingNotes
	}
	if product.Variety != nil {
		raw.VarietyRaw = *product.Variety
	}
	if product.Producer != nil {
		raw.ProducerRaw = *product.Producer
	}

	return &raw, nil
}
