package scraper

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the schema for a roaster's scraper configuration file.
// Config files live in roasters/configs/<slug>.yaml.
//
// The config-driven scraper supports the most common site patterns.
// For highly bespoke sites, implement the Scraper interface directly and
// register it in the registry.
type Config struct {
	// --- Roaster identity ---

	// Slug is the unique URL-safe identifier, e.g. "seven-seeds". Must match
	// the filename (without .yaml).
	Slug string `yaml:"slug"`

	// Name is the human-readable roaster name, e.g. "Seven Seeds".
	Name string `yaml:"name"`

	// State is the Australian state/territory abbreviation: VIC, NSW, QLD,
	// WA, SA, TAS, ACT, NT.
	State string `yaml:"state"`

	// Website is the roaster's homepage URL.
	Website string `yaml:"website"`

	// Active controls whether this roaster is included in scheduled scrape
	// runs. Set to false to pause without deleting the config.
	Active bool `yaml:"active"`

	// --- Scrape target ---

	// ShopURL is the URL of the coffee listing/shop page to scrape.
	// For paginated shops, this should be the first page.
	ShopURL string `yaml:"shop_url"`

	// PaginationType describes how subsequent pages are accessed.
	// Options: "none", "query_param", "next_button"
	PaginationType string `yaml:"pagination_type,omitempty"`

	// PaginationParam is the query parameter name used for page numbers,
	// e.g. "page". Only used when pagination_type is "query_param".
	PaginationParam string `yaml:"pagination_param,omitempty"`

	// MaxPages is a safety cap on how many pages to scrape. Defaults to 10.
	MaxPages int `yaml:"max_pages,omitempty"`

	// --- Extraction strategy ---

	// Strategy controls how product data is extracted.
	// Options:
	//   "jsonld"   — parse Schema.org Product JSON-LD (preferred)
	//   "selector" — use CSS selectors defined in Selectors below
	//   "shopify"  — use Shopify's products.json endpoint (auto-detected)
	Strategy string `yaml:"strategy"`

	// Selectors defines CSS selectors for strategy "selector".
	Selectors SelectorConfig `yaml:"selectors,omitempty"`

	// ShopifyHandle is only used for strategy "shopify". When set, only
	// products with this collection handle are fetched.
	ShopifyHandle string `yaml:"shopify_collection,omitempty"`

	// --- Field extraction hints ---

	// FieldHints provides additional hints for parsing specific fields when
	// the site uses non-standard patterns.
	FieldHints FieldHintConfig `yaml:"field_hints,omitempty"`

	// --- Ethical scraping ---

	// RateLimitSeconds overrides the default per-domain rate limit.
	// Must be >= 2. Defaults to 3.
	RateLimitSeconds int `yaml:"rate_limit_seconds,omitempty"`

	// Headers are additional HTTP request headers to send (e.g. Accept-Language).
	Headers map[string]string `yaml:"headers,omitempty"`
}

// SelectorConfig holds CSS selectors used with strategy "selector".
// All selectors are evaluated relative to each product card element.
type SelectorConfig struct {
	// ProductCard is the CSS selector for a single product listing container.
	// All other selectors below are evaluated within this element.
	// Example: ".product-item", "article.coffee-card"
	ProductCard string `yaml:"product_card"`

	// Name is the CSS selector for the coffee product name.
	// Example: "h2.product-title", ".product-name"
	Name string `yaml:"name"`

	// Price is the CSS selector for the price element.
	// The scraper will extract the text content and parse it.
	// Example: ".price", "span[data-price]"
	Price string `yaml:"price"`

	// Weight is the CSS selector for the weight/size option.
	// Example: "select.weight option:checked", ".variant-weight"
	Weight string `yaml:"weight,omitempty"`

	// Origin is the CSS selector for the origin text.
	// Example: ".product-meta .origin", "td.origin"
	Origin string `yaml:"origin,omitempty"`

	// Process is the CSS selector for the process method text.
	// Example: ".product-meta .process"
	Process string `yaml:"process,omitempty"`

	// RoastLevel is the CSS selector for roast level.
	// Example: ".roast-badge", ".product-meta .roast"
	RoastLevel string `yaml:"roast_level,omitempty"`

	// TastingNotes is the CSS selector for tasting notes text.
	// Example: ".tasting-notes", "p.notes"
	TastingNotes string `yaml:"tasting_notes,omitempty"`

	// ProductURL is the CSS selector for the product page link.
	// The scraper extracts the href attribute. If empty, falls back to the
	// closest ancestor <a> tag within the product card.
	// Example: "a.product-link", "h2 a"
	ProductURL string `yaml:"product_url,omitempty"`

	// ImageURL is the CSS selector for the product image.
	// The scraper extracts the src or data-src attribute.
	// Example: "img.product-image", ".product-card img"
	ImageURL string `yaml:"image_url,omitempty"`

	// InStock is a CSS selector that, if present in the DOM, indicates the
	// product is in stock. If this selector matches zero elements, the
	// product is considered out of stock.
	// Example: ".add-to-cart:not([disabled])", ".in-stock-badge"
	InStock string `yaml:"in_stock,omitempty"`
}

// FieldHintConfig provides parsing hints for non-standard field formats.
type FieldHintConfig struct {
	// PriceAttribute is the HTML attribute to read price from instead of
	// text content. Example: "data-price" (value in cents), "content"
	PriceAttribute string `yaml:"price_attribute,omitempty"`

	// PriceInCents indicates the price attribute value is already in cents.
	PriceInCents bool `yaml:"price_in_cents,omitempty"`

	// WeightInGrams indicates the weight value is already in grams (not
	// a string like "250g").
	WeightInGrams bool `yaml:"weight_in_grams,omitempty"`

	// OriginInTitle indicates that origin information is embedded in the
	// product name/title and should be parsed from there.
	OriginInTitle bool `yaml:"origin_in_title,omitempty"`

	// SkipPatterns is a list of product name substrings (case-insensitive)
	// that indicate the product should be skipped (e.g. merchandise, gift
	// cards, subscriptions).
	SkipPatterns []string `yaml:"skip_patterns,omitempty"`
}

// LoadConfig reads and validates a roaster config from a YAML file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	// Defaults
	if cfg.MaxPages == 0 {
		cfg.MaxPages = 10
	}
	if cfg.RateLimitSeconds == 0 {
		cfg.RateLimitSeconds = 3
	}
	if cfg.PaginationType == "" {
		cfg.PaginationType = "none"
	}

	return &cfg, nil
}

// LoadAllConfigs reads all *.yaml files in a directory and returns the
// parsed configs, skipping inactive ones unless includeInactive is true.
func LoadAllConfigs(dir string, includeInactive bool) ([]*Config, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob configs in %s: %w", dir, err)
	}

	var configs []*Config
	for _, path := range entries {
		// Skip the template file
		if filepath.Base(path) == "_template.yaml" {
			continue
		}

		cfg, err := LoadConfig(path)
		if err != nil {
			return nil, err
		}

		if !includeInactive && !cfg.Active {
			continue
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}

// Validate checks required fields and known enum values.
func (c *Config) Validate() error {
	if c.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.ShopURL == "" {
		return fmt.Errorf("shop_url is required")
	}

	validStates := map[string]bool{
		"VIC": true, "NSW": true, "QLD": true, "WA": true,
		"SA": true, "TAS": true, "ACT": true, "NT": true,
	}
	if c.State != "" && !validStates[c.State] {
		return fmt.Errorf("state %q is not a valid Australian state/territory", c.State)
	}

	validStrategies := map[string]bool{
		"jsonld": true, "selector": true, "shopify": true,
	}
	if !validStrategies[c.Strategy] {
		return fmt.Errorf("strategy must be one of: jsonld, selector, shopify (got %q)", c.Strategy)
	}

	if c.Strategy == "selector" && c.Selectors.ProductCard == "" {
		return fmt.Errorf("selectors.product_card is required when strategy is 'selector'")
	}

	validPagination := map[string]bool{
		"": true, "none": true, "query_param": true, "next_button": true,
	}
	if !validPagination[c.PaginationType] {
		return fmt.Errorf("pagination_type must be one of: none, query_param, next_button")
	}

	if c.RateLimitSeconds > 0 && c.RateLimitSeconds < 2 {
		return fmt.Errorf("rate_limit_seconds must be >= 2 (got %d)", c.RateLimitSeconds)
	}

	return nil
}
