package domain

// --- Auth ---

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required"`
	Password string `json:"password" binding:"required,min=8" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type UserResponse struct {
	ID    int32  `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type AuthResponse struct {
	User UserResponse `json:"user" validate:"required"`
}

// --- Scraper ---

// FetchMethod determines how a roaster's products are fetched.
type FetchMethod string

const (
	FetchShopifyJSON FetchMethod = "shopify_json"
	FetchHTML        FetchMethod = "html"
)

// RoasterConfig defines a roaster to scrape, loaded from roasters.yaml.
type RoasterConfig struct {
	Slug        string      `yaml:"slug"`
	Name        string      `yaml:"name"`
	Website     string      `yaml:"website"`
	ShopURL     string      `yaml:"shop_url"`
	State       string      `yaml:"state"`
	FetchMethod     FetchMethod `yaml:"fetch_method"`
	ProductType     string      `yaml:"product_type,omitempty"` // Shopify product_type filter
	ContentSelector string      `yaml:"content_selector,omitempty"` // CSS selector for product listing area (HTML fetch)
	Active          bool        `yaml:"active"`
}

// RoastersFile is the top-level structure of roasters.yaml.
type RoastersFile struct {
	Roasters []RoasterConfig `yaml:"roasters"`
}

// --- API Responses ---

type RoasterResponse struct {
	ID      int32  `json:"id"`
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Website string `json:"website"`
	State   string `json:"state,omitempty"`
}

type CoffeeResponse struct {
	ID           int64    `json:"id"`
	RoasterID    int32    `json:"roaster_id"`
	RoasterName  string   `json:"roaster_name,omitempty"`
	RoasterSlug  string   `json:"roaster_slug,omitempty"`
	Name         string   `json:"name"`
	ProductURL   string   `json:"product_url,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
	CountryCode  string   `json:"country_code,omitempty"`
	CountryName  string   `json:"country_name,omitempty"`
	RegionID     int32    `json:"region_id,omitempty"`
	RegionName   string   `json:"region_name,omitempty"`
	ProducerID   int32    `json:"producer_id,omitempty"`
	ProducerName string   `json:"producer_name,omitempty"`
	Process      string   `json:"process,omitempty"`
	RoastLevel   string   `json:"roast_level,omitempty"`
	TastingNotes []string `json:"tasting_notes,omitempty"`
	Variety      string   `json:"variety,omitempty"`
	Species      string   `json:"species,omitempty"`
	PriceCents   int32    `json:"price_cents,omitempty"`
	WeightGrams  int32    `json:"weight_grams,omitempty"`
	InStock      bool     `json:"in_stock"`
}

// CoffeeDetailResponse wraps a CoffeeResponse with similar coffees.
type CoffeeDetailResponse struct {
	CoffeeResponse
	SimilarCoffees []SimilarCoffee `json:"similar_coffees,omitempty"`
}

// SimilarCoffee is a lightweight coffee representation for the similar coffees section.
type SimilarCoffee struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	RoasterName  string   `json:"roaster_name"`
	RoasterSlug  string   `json:"roaster_slug"`
	ImageURL     string   `json:"image_url,omitempty"`
	Process      string   `json:"process,omitempty"`
	RoastLevel   string   `json:"roast_level,omitempty"`
	TastingNotes []string `json:"tasting_notes,omitempty"`
	Variety      string   `json:"variety,omitempty"`
	Score        float64  `json:"score"`
}

type CoffeeListResponse struct {
	Coffees    []CoffeeResponse `json:"coffees"`
	TotalCount int64            `json:"total_count"`
	Page       int32            `json:"page"`
	PageSize   int32            `json:"page_size"`
}

type RoasterListResponse struct {
	Roasters []RoasterResponse `json:"roasters"`
}

type RoasterDetailResponse struct {
	Roaster RoasterResponse  `json:"roaster"`
	Coffees []CoffeeResponse `json:"coffees"`
}

type StatsResponse struct {
	RoasterCount int64    `json:"roaster_count"`
	CoffeeCount  int64    `json:"coffee_count"`
	Origins      []string `json:"origins"`
}

// --- Countries, Regions, Producers ---

type CountryResponse struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	CoffeeCount int32  `json:"coffee_count"`
}

type CountryListResponse struct {
	Countries []CountryResponse `json:"countries"`
}

type CountryDetailResponse struct {
	Code    string           `json:"code"`
	Name    string           `json:"name"`
	Regions []RegionSummary  `json:"regions"`
	Coffees []CoffeeResponse `json:"coffees"`
}

type RegionSummary struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	CoffeeCount int32  `json:"coffee_count"`
}

type RegionDetailResponse struct {
	ID            int32            `json:"id"`
	Name          string           `json:"name"`
	CountryCode   string           `json:"country_code"`
	CountryName   string           `json:"country_name"`
	Latitude      *float64         `json:"latitude,omitempty"`
	Longitude     *float64         `json:"longitude,omitempty"`
	Coffees       []CoffeeResponse `json:"coffees"`
	NearbyRegions []NearbyRegion   `json:"nearby_regions,omitempty"`
}

type NearbyRegion struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	DistanceKm  int32  `json:"distance_km"`
	CoffeeCount int32  `json:"coffee_count"`
}

type ProducerResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	CoffeeCount int32  `json:"coffee_count"`
}

type ProducerListResponse struct {
	Producers []ProducerResponse `json:"producers"`
}

type ProducerDetailResponse struct {
	ID          int32            `json:"id"`
	Name        string           `json:"name"`
	CountryCode string           `json:"country_code,omitempty"`
	CountryName string           `json:"country_name,omitempty"`
	RegionName  string           `json:"region_name,omitempty"`
	Coffees     []CoffeeResponse `json:"coffees"`
}
