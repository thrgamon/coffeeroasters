package flavour

// Vector represents a coffee's position across six flavour dimensions.
// Values range from 0.0 (absent) to 1.0 (dominant).
type Vector struct {
	Fruity float64
	Sweet  float64
	Bright float64
	Body   float64
	Floral float64
	Earthy float64
}

// ProcessWeights returns the flavour vector for a canonical process value.
func ProcessWeights(process string) Vector {
	v, ok := processWeights[process]
	if !ok {
		return Vector{}
	}
	return v
}

// RoastWeights returns the flavour vector for a canonical roast level.
func RoastWeights(roast string) Vector {
	v, ok := roastWeights[roast]
	if !ok {
		return Vector{}
	}
	return v
}

// OriginWeights returns the flavour vector for a country code.
func OriginWeights(countryCode string) Vector {
	v, ok := originWeights[countryCode]
	if !ok {
		return Vector{}
	}
	return v
}

// VarietyWeights returns the flavour vector for a canonical variety.
func VarietyWeights(variety string) Vector {
	v, ok := varietyWeights[variety]
	if !ok {
		return Vector{}
	}
	return v
}

var processWeights = map[string]Vector{
	"washed":       {Fruity: 0.3, Sweet: 0.3, Bright: 0.8, Body: 0.3, Floral: 0.7, Earthy: 0.1},
	"natural":      {Fruity: 0.9, Sweet: 0.6, Bright: 0.4, Body: 0.8, Floral: 0.3, Earthy: 0.1},
	"honey":        {Fruity: 0.6, Sweet: 0.7, Bright: 0.5, Body: 0.6, Floral: 0.4, Earthy: 0.1},
	"anaerobic":    {Fruity: 0.8, Sweet: 0.5, Bright: 0.5, Body: 0.7, Floral: 0.4, Earthy: 0.2},
	"wet-hulled":   {Fruity: 0.1, Sweet: 0.3, Bright: 0.1, Body: 0.9, Floral: 0.0, Earthy: 0.9},
	"experimental": {Fruity: 0.7, Sweet: 0.5, Bright: 0.5, Body: 0.6, Floral: 0.4, Earthy: 0.2},
}

var roastWeights = map[string]Vector{
	"light":        {Fruity: 0.7, Sweet: 0.3, Bright: 0.9, Body: 0.2, Floral: 0.8, Earthy: 0.0},
	"medium-light": {Fruity: 0.5, Sweet: 0.5, Bright: 0.7, Body: 0.4, Floral: 0.6, Earthy: 0.1},
	"medium":       {Fruity: 0.3, Sweet: 0.7, Bright: 0.5, Body: 0.6, Floral: 0.3, Earthy: 0.2},
	"medium-dark":  {Fruity: 0.2, Sweet: 0.8, Bright: 0.2, Body: 0.8, Floral: 0.1, Earthy: 0.4},
	"dark":         {Fruity: 0.1, Sweet: 0.6, Bright: 0.1, Body: 0.9, Floral: 0.0, Earthy: 0.6},
}

var originWeights = map[string]Vector{
	// Africa
	"ET": {Fruity: 0.7, Sweet: 0.4, Bright: 0.8, Body: 0.3, Floral: 0.9, Earthy: 0.0},
	"KE": {Fruity: 0.8, Sweet: 0.3, Bright: 0.9, Body: 0.6, Floral: 0.5, Earthy: 0.1},
	"RW": {Fruity: 0.6, Sweet: 0.5, Bright: 0.7, Body: 0.5, Floral: 0.6, Earthy: 0.0},
	"BI": {Fruity: 0.7, Sweet: 0.4, Bright: 0.7, Body: 0.4, Floral: 0.5, Earthy: 0.0},
	"TZ": {Fruity: 0.5, Sweet: 0.5, Bright: 0.6, Body: 0.5, Floral: 0.4, Earthy: 0.1},
	"UG": {Fruity: 0.5, Sweet: 0.5, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.2},
	"CD": {Fruity: 0.6, Sweet: 0.4, Bright: 0.6, Body: 0.5, Floral: 0.5, Earthy: 0.1},
	"MW": {Fruity: 0.5, Sweet: 0.5, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	// Americas
	"CO": {Fruity: 0.5, Sweet: 0.7, Bright: 0.6, Body: 0.6, Floral: 0.3, Earthy: 0.1},
	"BR": {Fruity: 0.2, Sweet: 0.8, Bright: 0.2, Body: 0.7, Floral: 0.1, Earthy: 0.3},
	"GT": {Fruity: 0.4, Sweet: 0.6, Bright: 0.6, Body: 0.7, Floral: 0.2, Earthy: 0.2},
	"CR": {Fruity: 0.5, Sweet: 0.6, Bright: 0.6, Body: 0.4, Floral: 0.3, Earthy: 0.1},
	"PA": {Fruity: 0.6, Sweet: 0.6, Bright: 0.6, Body: 0.4, Floral: 0.5, Earthy: 0.0},
	"HN": {Fruity: 0.4, Sweet: 0.6, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"SV": {Fruity: 0.4, Sweet: 0.6, Bright: 0.6, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"NI": {Fruity: 0.3, Sweet: 0.6, Bright: 0.4, Body: 0.5, Floral: 0.2, Earthy: 0.2},
	"MX": {Fruity: 0.3, Sweet: 0.6, Bright: 0.4, Body: 0.5, Floral: 0.2, Earthy: 0.2},
	"PE": {Fruity: 0.4, Sweet: 0.6, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"EC": {Fruity: 0.5, Sweet: 0.5, Bright: 0.5, Body: 0.5, Floral: 0.4, Earthy: 0.1},
	"BO": {Fruity: 0.4, Sweet: 0.5, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	// Asia-Pacific
	"ID": {Fruity: 0.1, Sweet: 0.3, Bright: 0.1, Body: 0.9, Floral: 0.0, Earthy: 0.9},
	"PG": {Fruity: 0.3, Sweet: 0.5, Bright: 0.4, Body: 0.6, Floral: 0.2, Earthy: 0.4},
	"IN": {Fruity: 0.2, Sweet: 0.4, Bright: 0.2, Body: 0.7, Floral: 0.1, Earthy: 0.6},
	"MM": {Fruity: 0.4, Sweet: 0.5, Bright: 0.4, Body: 0.5, Floral: 0.3, Earthy: 0.3},
	"CN": {Fruity: 0.3, Sweet: 0.5, Bright: 0.4, Body: 0.5, Floral: 0.3, Earthy: 0.3},
	// Middle East
	"YE": {Fruity: 0.6, Sweet: 0.5, Bright: 0.5, Body: 0.6, Floral: 0.4, Earthy: 0.3},
}

var varietyWeights = map[string]Vector{
	"bourbon":        {Fruity: 0.5, Sweet: 0.7, Bright: 0.5, Body: 0.6, Floral: 0.3, Earthy: 0.1},
	"typica":         {Fruity: 0.4, Sweet: 0.6, Bright: 0.5, Body: 0.5, Floral: 0.4, Earthy: 0.1},
	"caturra":        {Fruity: 0.5, Sweet: 0.5, Bright: 0.7, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"catuai":         {Fruity: 0.4, Sweet: 0.6, Bright: 0.5, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"gesha":          {Fruity: 0.6, Sweet: 0.5, Bright: 0.8, Body: 0.2, Floral: 0.9, Earthy: 0.0},
	"sl28":           {Fruity: 0.8, Sweet: 0.5, Bright: 0.9, Body: 0.6, Floral: 0.4, Earthy: 0.0},
	"sl34":           {Fruity: 0.7, Sweet: 0.5, Bright: 0.8, Body: 0.7, Floral: 0.4, Earthy: 0.0},
	"pacamara":       {Fruity: 0.7, Sweet: 0.5, Bright: 0.6, Body: 0.6, Floral: 0.5, Earthy: 0.1},
	"heirloom":       {Fruity: 0.7, Sweet: 0.4, Bright: 0.7, Body: 0.3, Floral: 0.8, Earthy: 0.0},
	"castillo":       {Fruity: 0.5, Sweet: 0.5, Bright: 0.6, Body: 0.5, Floral: 0.3, Earthy: 0.1},
	"catimor":        {Fruity: 0.3, Sweet: 0.5, Bright: 0.4, Body: 0.6, Floral: 0.2, Earthy: 0.2},
	"mundo-novo":     {Fruity: 0.3, Sweet: 0.6, Bright: 0.4, Body: 0.6, Floral: 0.2, Earthy: 0.1},
	"yellow-bourbon": {Fruity: 0.5, Sweet: 0.7, Bright: 0.5, Body: 0.5, Floral: 0.4, Earthy: 0.1},
	"pink-bourbon":   {Fruity: 0.6, Sweet: 0.6, Bright: 0.6, Body: 0.5, Floral: 0.5, Earthy: 0.0},
	"red-bourbon":    {Fruity: 0.5, Sweet: 0.7, Bright: 0.5, Body: 0.6, Floral: 0.3, Earthy: 0.1},
	"sidra":          {Fruity: 0.7, Sweet: 0.5, Bright: 0.6, Body: 0.5, Floral: 0.6, Earthy: 0.0},
	"wush-wush":      {Fruity: 0.6, Sweet: 0.4, Bright: 0.7, Body: 0.3, Floral: 0.7, Earthy: 0.0},
}
