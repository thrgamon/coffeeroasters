package scraper

import "strings"

// excludedKeywords is a safety-net filter applied after LLM extraction.
// Products whose name contains any of these substrings (case-insensitive) are
// discarded even if the LLM marked them as coffee.
var excludedKeywords = []string{
	"drip bag",
	"drip coffee",
	"instant coffee",
	"cold brew",
	"cold wolff",
	"concentrate",
	"ready to drink",
	"capsule",
	"pod",
	"grinder",
	"equipment",
	"gift card",
	"voucher",
	"merch",
	"subscription",
	"sample pack",
	"bundle",
	"rtd",
	"parachute",
}

// isExcludedByKeyword returns true if the product name matches any excluded keyword.
func isExcludedByKeyword(name string) bool {
	lower := strings.ToLower(name)
	for _, kw := range excludedKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
