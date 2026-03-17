package normalise

import "strings"

// Variety is a canonical coffee cultivar name.
type Variety string

const (
	VarietyBourbon      Variety = "bourbon"
	VarietyTypica       Variety = "typica"
	VarietyCaturra      Variety = "caturra"
	VarietyCatuai       Variety = "catuai"
	VarietyGesha        Variety = "gesha"
	VarietySL28         Variety = "sl28"
	VarietySL34         Variety = "sl34"
	VarietyPacamara     Variety = "pacamara"
	VarietyMaragogipe   Variety = "maragogipe"
	VarietyHeirloom     Variety = "heirloom"
	VarietyCastillo     Variety = "castillo"
	VarietyColombia     Variety = "colombia"
	VarietyJava         Variety = "java"
	VarietyRuiru11      Variety = "ruiru11"
	VarietyBatian       Variety = "batian"
	VarietyCatimor      Variety = "catimor"
	VarietyMarsellesa   Variety = "marsellesa"
	VarietyParainema    Variety = "parainema"
	VarietyObata        Variety = "obata"
	VarietyMundoNovo    Variety = "mundo-novo"
	VarietyYellowBurbon Variety = "yellow-bourbon"
	VarietyPinkBourbon  Variety = "pink-bourbon"
	VarietyRedBourbon   Variety = "red-bourbon"
	VarietyTabi         Variety = "tabi"
	VarietySidra        Variety = "sidra"
	VarietyWushWush     Variety = "wush-wush"
	Variety74110        Variety = "74110"
	Variety74112        Variety = "74112"
	Variety74158        Variety = "74158"
	VarietyPacas        Variety = "pacas"
	VarietyMaracaturra  Variety = "maracaturra"
	VarietyUnknown      Variety = ""
)

// Species is a canonical coffee species.
type Species string

const (
	SpeciesArabica  Species = "arabica"
	SpeciesRobusta  Species = "robusta"
	SpeciesLiberica Species = "liberica"
	SpeciesUnknown  Species = ""
)

// varietyAliases maps common spellings and abbreviations to canonical variety
// values. Keys must be lowercased.
var varietyAliases = map[string]Variety{
	// Gesha/Geisha
	"gesha":   VarietyGesha,
	"geisha":  VarietyGesha,
	"gescha":  VarietyGesha,

	// Bourbon variants
	"bourbon":       VarietyBourbon,
	"borbon":        VarietyBourbon,
	"bourbon rouge": VarietyRedBourbon,
	"red bourbon":   VarietyRedBourbon,
	"yellow bourbon": VarietyYellowBurbon,
	"bourbon amarillo": VarietyYellowBurbon,
	"pink bourbon":  VarietyPinkBourbon,

	// Typica
	"typica":  VarietyTypica,
	"tipica":  VarietyTypica,

	// Caturra
	"caturra":      VarietyCaturra,
	"caturra rojo": VarietyCaturra,

	// Catuai
	"catuai":        VarietyCatuai,
	"catuai rojo":   VarietyCatuai,
	"catuai amarillo": VarietyCatuai,

	// SL28 / SL34
	"sl28":  VarietySL28,
	"sl-28": VarietySL28,
	"sl 28": VarietySL28,
	"sl34":  VarietySL34,
	"sl-34": VarietySL34,
	"sl 34": VarietySL34,

	// Pacamara
	"pacamara": VarietyPacamara,

	// Maragogipe
	"maragogipe":  VarietyMaragogipe,
	"maragogype":  VarietyMaragogipe,
	"maracaturra": VarietyMaracaturra,

	// Heirloom
	"heirloom":           VarietyHeirloom,
	"heirloom varieties": VarietyHeirloom,
	"ethiopian heirloom": VarietyHeirloom,
	"landraces":          VarietyHeirloom,

	// Castillo
	"castillo": VarietyCastillo,

	// Colombia
	"colombia": VarietyColombia,

	// Java
	"java": VarietyJava,

	// Ruiru 11
	"ruiru 11": VarietyRuiru11,
	"ruiru11":  VarietyRuiru11,

	// Batian
	"batian": VarietyBatian,

	// Catimor
	"catimor": VarietyCatimor,

	// Marsellesa
	"marsellesa": VarietyMarsellesa,

	// Parainema
	"parainema": VarietyParainema,

	// Obata
	"obata": VarietyObata,

	// Mundo Novo
	"mundo novo":  VarietyMundoNovo,
	"mundo-novo":  VarietyMundoNovo,

	// Tabi
	"tabi": VarietyTabi,

	// Sidra
	"sidra": VarietySidra,

	// Wush Wush
	"wush wush":  VarietyWushWush,
	"wush-wush":  VarietyWushWush,

	// Ethiopian selections
	"74110": Variety74110,
	"74112": Variety74112,
	"74158": Variety74158,

	// Pacas
	"pacas": VarietyPacas,
}

// varietySpecies maps known varieties to their species. Most specialty
// varieties are arabica.
var varietySpecies = map[Variety]Species{
	VarietyBourbon:      SpeciesArabica,
	VarietyTypica:       SpeciesArabica,
	VarietyCaturra:      SpeciesArabica,
	VarietyCatuai:       SpeciesArabica,
	VarietyGesha:        SpeciesArabica,
	VarietySL28:         SpeciesArabica,
	VarietySL34:         SpeciesArabica,
	VarietyPacamara:     SpeciesArabica,
	VarietyMaragogipe:   SpeciesArabica,
	VarietyHeirloom:     SpeciesArabica,
	VarietyCastillo:     SpeciesArabica,
	VarietyColombia:     SpeciesArabica,
	VarietyJava:         SpeciesArabica,
	VarietyRuiru11:      SpeciesArabica,
	VarietyBatian:       SpeciesArabica,
	VarietyMarsellesa:   SpeciesArabica,
	VarietyParainema:    SpeciesArabica,
	VarietyObata:        SpeciesArabica,
	VarietyMundoNovo:    SpeciesArabica,
	VarietyYellowBurbon: SpeciesArabica,
	VarietyPinkBourbon:  SpeciesArabica,
	VarietyRedBourbon:   SpeciesArabica,
	VarietyTabi:         SpeciesArabica,
	VarietySidra:        SpeciesArabica,
	VarietyWushWush:     SpeciesArabica,
	Variety74110:        SpeciesArabica,
	Variety74112:        SpeciesArabica,
	Variety74158:        SpeciesArabica,
	VarietyPacas:        SpeciesArabica,
	VarietyMaracaturra:  SpeciesArabica,
	VarietyCatimor:      SpeciesArabica,
}

// robusta aliases that map to robusta species
var robustaAliases = map[string]bool{
	"robusta":           true,
	"conilon":           true,
	"coffea canephora":  true,
}

// NormaliseVariety normalises a raw variety string into canonical variety and
// species values. Multi-variety blends (comma or slash separated) are stored
// as comma-joined canonical names. Species is inferred from the first
// recognised variety.
func NormaliseVariety(raw string) (variety string, species string) {
	key := normaliseKey(raw)
	if key == "" {
		return "", ""
	}

	// Split on comma/slash for multi-variety
	separators := strings.NewReplacer("/", ",", " & ", ",", " and ", ",")
	key = separators.Replace(key)

	parts := strings.Split(key, ",")
	var canonical []string
	var firstSpecies Species

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		v := lookupVariety(p)
		if v == VarietyUnknown {
			// No match: return empty so LLM classifier can handle it
			return "", ""
		}

		canonical = append(canonical, string(v))

		if firstSpecies == SpeciesUnknown {
			firstSpecies = inferSpecies(v, p)
		}
	}

	if len(canonical) == 0 {
		return "", ""
	}

	return strings.Join(canonical, ","), string(firstSpecies)
}

func lookupVariety(key string) Variety {
	// Exact match
	if v, ok := varietyAliases[key]; ok {
		return v
	}

	// Partial match fallback
	for alias, v := range varietyAliases {
		if strings.Contains(key, alias) {
			return v
		}
	}

	return VarietyUnknown
}

func inferSpecies(v Variety, raw string) Species {
	// Check robusta aliases first (for raw input like "robusta", "conilon")
	if robustaAliases[raw] {
		return SpeciesRobusta
	}

	if s, ok := varietySpecies[v]; ok {
		return s
	}

	return SpeciesUnknown
}
