// Package normalise converts raw scraped strings into canonical, queryable
// values. All normalisation is deterministic and testable — no network calls.
//
// The normalisation pipeline is:
//  1. Raw string arrives from scraper
//  2. Lowercased, trimmed, common punctuation removed
//  3. Matched against a reference lookup (YAML file or embedded map)
//  4. Canonical value stored alongside the original raw string in the DB
package normalise

import (
	"strings"
)

// Process values — canonical enum stored in DB.
const (
	ProcessWashed       = "washed"
	ProcessNatural      = "natural"
	ProcessHoney        = "honey"
	ProcessAnaerobic    = "anaerobic"
	ProcessWetHulled    = "wet-hulled"
	ProcessExperimental = "experimental"
	ProcessUnknown      = ""
)

// RoastLevel values — canonical enum stored in DB.
const (
	RoastLight      = "light"
	RoastMedLight   = "medium-light"
	RoastMedium     = "medium"
	RoastMedDark    = "medium-dark"
	RoastDark       = "dark"
	RoastUnknown    = ""
)

// processAliases maps raw strings (lowercased) to canonical process values.
// Extend this as new variants are encountered in the wild.
var processAliases = map[string]string{
	// Washed / wet process
	"washed":         ProcessWashed,
	"wet process":    ProcessWashed,
	"wet-process":    ProcessWashed,
	"fully washed":   ProcessWashed,
	"full washed":    ProcessWashed,
	"wet":            ProcessWashed,

	// Natural / dry process
	"natural":        ProcessNatural,
	"naturale":       ProcessNatural,
	"dry process":    ProcessNatural,
	"dry-process":    ProcessNatural,
	"dried":          ProcessNatural,
	"sun dried":      ProcessNatural,
	"sun-dried":      ProcessNatural,

	// Honey / pulped natural
	"honey":          ProcessHoney,
	"pulped natural": ProcessHoney,
	"pulped":         ProcessHoney,
	"yellow honey":   ProcessHoney,
	"red honey":      ProcessHoney,
	"black honey":    ProcessHoney,
	"white honey":    ProcessHoney,

	// Anaerobic
	"anaerobic":              ProcessAnaerobic,
	"anaerobic natural":      ProcessAnaerobic,
	"anaerobic washed":       ProcessAnaerobic,
	"anaerobic fermentation": ProcessAnaerobic,
	"extended fermentation":  ProcessAnaerobic,

	// Wet-hulled (Giling Basah)
	"wet hulled":   ProcessWetHulled,
	"wet-hulled":   ProcessWetHulled,
	"giling basah": ProcessWetHulled,

	// Experimental
	"experimental":    ProcessExperimental,
	"carbonic":        ProcessExperimental,
	"carbonic maceration": ProcessExperimental,
	"lactic":          ProcessExperimental,
}

// roastAliases maps raw strings (lowercased) to canonical roast level values.
var roastAliases = map[string]string{
	// Light
	"light":         RoastLight,
	"light roast":   RoastLight,
	"filter":        RoastLight,
	"filter roast":  RoastLight,
	"omni":          RoastLight,
	"omni roast":    RoastLight,

	// Medium-light
	"medium light":  RoastMedLight,
	"medium-light":  RoastMedLight,
	"med-light":     RoastMedLight,
	"espresso light": RoastMedLight,

	// Medium
	"medium":        RoastMedium,
	"medium roast":  RoastMedium,
	"med":           RoastMedium,
	"espresso":      RoastMedium,
	"espresso roast": RoastMedium,
	"all rounder":   RoastMedium,
	"all-rounder":   RoastMedium,

	// Medium-dark
	"medium dark":   RoastMedDark,
	"medium-dark":   RoastMedDark,
	"med-dark":      RoastMedDark,
	"full city":     RoastMedDark,

	// Dark
	"dark":          RoastDark,
	"dark roast":    RoastDark,
	"french roast":  RoastDark,
	"italian roast": RoastDark,
}

// NormaliseProcess maps a raw process string to a canonical value.
// Returns ProcessUnknown ("") if no match is found.
func NormaliseProcess(raw string) string {
	key := normaliseKey(raw)
	if v, ok := processAliases[key]; ok {
		return v
	}
	// Partial match fallback
	for alias, canonical := range processAliases {
		if strings.Contains(key, alias) {
			return canonical
		}
	}
	return ProcessUnknown
}

// NormaliseRoastLevel maps a raw roast level string to a canonical value.
// Returns RoastUnknown ("") if no match is found.
func NormaliseRoastLevel(raw string) string {
	key := normaliseKey(raw)
	if v, ok := roastAliases[key]; ok {
		return v
	}
	for alias, canonical := range roastAliases {
		if strings.Contains(key, alias) {
			return canonical
		}
	}
	return RoastUnknown
}

// NormaliseTastingNotes splits a raw tasting notes string into individual
// clean note tokens. Input may be comma-separated, slash-separated, or
// use "and"/"&" conjunctions.
//
// Example: "blueberry, dark choc & jasmine" → ["blueberry", "dark chocolate", "jasmine"]
func NormaliseTastingNotes(raw string) []string {
	if raw == "" {
		return nil
	}

	// Replace common separators with comma
	r := raw
	r = strings.ReplaceAll(r, " & ", ", ")
	r = strings.ReplaceAll(r, " and ", ", ")
	r = strings.ReplaceAll(r, " / ", ", ")
	r = strings.ReplaceAll(r, "/", ", ")
	r = strings.ReplaceAll(r, "|", ", ")

	parts := strings.Split(r, ",")
	var notes []string
	for _, p := range parts {
		note := strings.TrimSpace(p)
		if note != "" {
			notes = append(notes, applyNoteAliases(note))
		}
	}
	return notes
}

// NormalisePriceAUD parses a raw price string and returns cents (int64).
// Returns 0, false if the string cannot be parsed.
//
// Handles: "$32.00", "32.00", "32", "$32", "AUD 32.00"
func NormalisePriceAUD(raw string) (cents int64, ok bool) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimPrefix(s, "AUD")
	s = strings.TrimPrefix(s, "AU$")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")

	if s == "" {
		return 0, false
	}

	// Parse as float to handle "$32.50"
	var dollars float64
	_, err := parseFloat(s, &dollars)
	if err != nil {
		return 0, false
	}

	return int64(dollars * 100), true
}

// NormaliseWeightGrams parses a raw weight string and returns grams (int).
// Returns 0, false if unparseable.
//
// Handles: "250g", "1kg", "250 g", "1 kg", "250"
func NormaliseWeightGrams(raw string) (grams int, ok bool) {
	s := strings.ToLower(strings.TrimSpace(raw))

	if strings.HasSuffix(s, "kg") {
		s = strings.TrimSuffix(s, "kg")
		s = strings.TrimSpace(s)
		var kg float64
		if _, err := parseFloat(s, &kg); err != nil {
			return 0, false
		}
		return int(kg * 1000), true
	}

	s = strings.TrimSuffix(s, "g")
	s = strings.TrimSpace(s)
	var g float64
	if _, err := parseFloat(s, &g); err != nil {
		return 0, false
	}
	return int(g), true
}

// normaliseKey lowercases and trims a string for lookup.
func normaliseKey(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// noteAliases maps informal tasting note terms to cleaner display names.
var noteAliases = map[string]string{
	"choc":        "chocolate",
	"dark choc":   "dark chocolate",
	"milk choc":   "milk chocolate",
	"citrus":      "citrus",
	"berry":       "berry",
	"floral":      "floral",
	"nuts":        "nutty",
	"nut":         "nutty",
	"caramel":     "caramel",
	"stone fruit": "stone fruit",
}

func applyNoteAliases(note string) string {
	key := normaliseKey(note)
	if v, ok := noteAliases[key]; ok {
		return v
	}
	return note
}

// parseFloat is a minimal float parser to avoid importing strconv in the hot path.
func parseFloat(s string, out *float64) (n int, err error) {
	if len(s) == 0 {
		return 0, &parseError{s}
	}
	var integer, decimal int64
	var hasDecimal bool
	var decimalPlaces int64 = 1

	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if hasDecimal {
				return i, &parseError{s}
			}
			hasDecimal = true
			continue
		}
		if c < '0' || c > '9' {
			return i, &parseError{s}
		}
		if hasDecimal {
			decimal = decimal*10 + int64(c-'0')
			decimalPlaces *= 10
		} else {
			integer = integer*10 + int64(c-'0')
		}
	}

	*out = float64(integer) + float64(decimal)/float64(decimalPlaces)
	return i, nil
}

type parseError struct{ s string }

func (e *parseError) Error() string { return "cannot parse float: " + e.s }
