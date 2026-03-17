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

// countryAliases maps lowercased origin strings to ISO 3166-1 alpha-2 codes.
var countryAliases = map[string]string{
	// Ethiopia
	"ethiopia": "ET", "ethiopian": "ET", "ethiopa": "ET",
	// Colombia
	"colombia": "CO", "colombian": "CO", "columbian": "CO", "columbia": "CO",
	// Brazil
	"brazil": "BR", "brazilian": "BR", "brasil": "BR",
	// Kenya
	"kenya": "KE", "kenyan": "KE",
	// Guatemala
	"guatemala": "GT", "guatemalan": "GT",
	// Costa Rica
	"costa rica": "CR", "costa rican": "CR",
	// Panama
	"panama": "PA", "panamanian": "PA",
	// Honduras
	"honduras": "HN", "honduran": "HN",
	// El Salvador
	"el salvador": "SV", "salvadoran": "SV",
	// Nicaragua
	"nicaragua": "NI", "nicaraguan": "NI",
	// Mexico
	"mexico": "MX", "mexican": "MX",
	// Peru
	"peru": "PE", "peruvian": "PE",
	// Bolivia
	"bolivia": "BO", "bolivian": "BO",
	// Ecuador
	"ecuador": "EC", "ecuadorian": "EC",
	// Rwanda
	"rwanda": "RW", "rwandan": "RW",
	// Burundi
	"burundi": "BI", "burundian": "BI",
	// Tanzania
	"tanzania": "TZ", "tanzanian": "TZ",
	// Uganda
	"uganda": "UG", "ugandan": "UG",
	// DR Congo
	"dr congo": "CD", "congo": "CD", "drc": "CD",
	// Malawi
	"malawi": "MW", "malawian": "MW",
	// Zambia
	"zambia": "ZM", "zambian": "ZM",
	// Indonesia
	"indonesia": "ID", "indonesian": "ID",
	// Papua New Guinea
	"papua new guinea": "PG", "png": "PG",
	// India
	"india": "IN", "indian": "IN",
	// Myanmar
	"myanmar": "MM", "burma": "MM",
	// Vietnam
	"vietnam": "VN", "vietnamese": "VN",
	// China
	"china": "CN", "chinese": "CN",
	// Thailand
	"thailand": "TH", "thai": "TH",
	// Laos
	"laos": "LA", "lao": "LA",
	// Philippines
	"philippines": "PH", "philippine": "PH",
	// Taiwan
	"taiwan": "TW",
	// Yemen
	"yemen": "YE", "yemeni": "YE",
	// Jamaica
	"jamaica": "JM", "jamaican": "JM",
	// Haiti
	"haiti": "HT", "haitian": "HT",
	// Dominican Republic
	"dominican republic": "DO",
	// Cuba
	"cuba": "CU", "cuban": "CU",
}

// regionCountryMap maps well-known region names to their country code. Used as
// a fallback when the origin string is just a region name without a country.
var regionCountryMap = map[string]string{
	"yirgacheffe": "ET", "guji": "ET", "sidamo": "ET", "sidama": "ET",
	"gedeo": "ET", "gedeo zone": "ET", "limu": "ET", "jimma": "ET",
	"harrar": "ET", "harar": "ET",
	"huila": "CO", "nariño": "CO", "narino": "CO", "tolima": "CO",
	"cauca": "CO", "antioquia": "CO", "quindio": "CO",
	"cerrado": "BR", "mogiana": "BR", "sul de minas": "BR",
	"nyeri": "KE", "kirinyaga": "KE", "kiambu": "KE", "muranga": "KE",
	"antigua": "GT", "huehuetenango": "GT", "acatenango": "GT",
	"tarrazú": "CR", "tarrazu": "CR", "west valley": "CR",
	"boquete": "PA", "geisha": "PA",
	"copan": "HN", "santa barbara": "HN", "comayagua": "HN",
	"apaneca": "SV",
	"jinotega": "NI", "matagalpa": "NI",
	"oaxaca": "MX", "chiapas": "MX", "veracruz": "MX",
	"cajamarca": "PE", "san martin": "PE",
	"caranavi": "BO", "yungas": "BO",
	"huye": "RW", "nyamasheke": "RW",
	"kayanza": "BI", "ngozi": "BI",
	"mbeya": "TZ", "kilimanjaro": "TZ", "arusha": "TZ",
	"mount elgon": "UG", "bugisu": "UG", "rwenzori": "UG",
	"kivu": "CD",
	"sumatra": "ID", "java": "ID", "sulawesi": "ID",
	"aceh": "ID", "gayo": "ID", "lintong": "ID", "toraja": "ID",
	"yunnan": "CN",
	"eastern highlands": "PG",
}

// NormaliseOrigin maps a raw origin string and optional region string to an
// ISO 3166-1 alpha-2 country code and a cleaned region name.
//
// Strategy:
//  1. Try exact match of full originRaw against countryAliases
//  2. Try prefix match (e.g. "Ethiopia Yirgacheffe" matches "ethiopia")
//  3. Try regionCountryMap for well-known region names
//  4. If regionRaw is non-empty, pass it through as the region name
//  5. If country was matched via prefix, remainder is used as region
func NormaliseOrigin(originRaw, regionRaw string) (countryCode, regionName string) {
	origin := normaliseKey(originRaw)
	if origin == "" {
		return "", ""
	}

	// 1. Exact match on full string
	if code, ok := countryAliases[origin]; ok {
		countryCode = code
	}

	// 2. Prefix match: longest match wins
	if countryCode == "" {
		var bestLen int
		for alias, code := range countryAliases {
			if len(alias) > bestLen && strings.HasPrefix(origin, alias) {
				// Ensure the alias ends at a word boundary
				rest := origin[len(alias):]
				if rest == "" || rest[0] == ' ' || rest[0] == ',' || rest[0] == '-' {
					countryCode = code
					bestLen = len(alias)
				}
			}
		}
		// Extract region from remainder
		if bestLen > 0 && regionRaw == "" {
			rest := strings.TrimSpace(origin[bestLen:])
			rest = strings.TrimLeft(rest, ",- ")
			if rest != "" {
				regionName = rest
			}
		}
	}

	// 3. Try regionCountryMap
	if countryCode == "" {
		if code, ok := regionCountryMap[origin]; ok {
			countryCode = code
			if regionRaw == "" {
				regionName = origin
			}
		}
	}

	// regionRaw takes precedence over extracted region
	if regionRaw != "" {
		regionName = strings.TrimSpace(regionRaw)
	}

	return countryCode, regionName
}

// normaliseKey lowercases and trims a string for lookup.
func normaliseKey(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// noteAliases maps informal tasting note terms to cleaner display names.
// Aligned with SCA Flavour Wheel tier-2/tier-3 descriptors commonly seen
// on roaster websites.
var noteAliases = map[string]string{
	// Chocolate
	"choc":            "chocolate",
	"dark choc":       "dark chocolate",
	"milk choc":       "milk chocolate",
	"cocoa":           "chocolate",
	"cacao":           "chocolate",
	"dark chocolate":  "dark chocolate",
	"milk chocolate":  "milk chocolate",
	"white chocolate": "white chocolate",

	// Fruit
	"citrus":        "citrus",
	"citric":        "citrus",
	"berry":         "berry",
	"berries":       "berry",
	"stonefruit":    "stone fruit",
	"stone fruit":   "stone fruit",
	"tropical":      "tropical fruit",
	"tropical fruit": "tropical fruit",
	"dried fruit":   "dried fruit",

	// Nutty
	"nuts":     "nutty",
	"nut":      "nutty",
	"almond":   "almond",
	"hazelnut": "hazelnut",
	"peanut":   "peanut",
	"walnut":   "walnut",

	// Sweet
	"caramel":     "caramel",
	"brown sugar": "brown sugar",
	"molasses":    "molasses",
	"maple":       "maple syrup",
	"panela":      "panela",
	"toffee":      "toffee",

	// Floral
	"floral":      "floral",
	"jasmine":     "jasmine",
	"rose":        "rose",
	"lavender":    "lavender",
	"elderflower": "elderflower",

	// Spice
	"cinnamon":     "cinnamon",
	"clove":        "clove",
	"cardamom":     "cardamom",
	"black pepper": "black pepper",

	// Other
	"wine":     "winey",
	"boozy":    "winey",
	"tea-like": "tea-like",
	"herbal":   "herbal",
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
