package similarity

// flavourGroups maps SCA tier-3 descriptors to their tier-2 parent group.
// Used by the similarity scorer to give partial credit when two coffees share
// the same flavour family but different specific notes.
var flavourGroups = map[string]string{
	// Berry
	"blueberry":   "berry",
	"strawberry":  "berry",
	"raspberry":   "berry",
	"blackberry":  "berry",
	"cranberry":   "berry",

	// Citrus
	"lemon":       "citrus",
	"lime":        "citrus",
	"grapefruit":  "citrus",
	"orange":      "citrus",
	"tangerine":   "citrus",
	"mandarin":    "citrus",
	"bergamot":    "citrus",

	// Stone fruit
	"peach":       "stone fruit",
	"apricot":     "stone fruit",
	"plum":        "stone fruit",
	"cherry":      "stone fruit",
	"nectarine":   "stone fruit",

	// Tropical fruit
	"mango":       "tropical fruit",
	"pineapple":   "tropical fruit",
	"papaya":      "tropical fruit",
	"passionfruit": "tropical fruit",
	"lychee":      "tropical fruit",
	"guava":       "tropical fruit",
	"coconut":     "tropical fruit",

	// Dried fruit
	"raisin":      "dried fruit",
	"fig":         "dried fruit",
	"date":        "dried fruit",
	"prune":       "dried fruit",

	// Chocolate
	"dark chocolate": "chocolate",
	"milk chocolate": "chocolate",
	"white chocolate": "chocolate",
	"cocoa":          "chocolate",
	"cacao":          "chocolate",

	// Nutty
	"almond":   "nutty",
	"hazelnut": "nutty",
	"peanut":   "nutty",
	"walnut":   "nutty",
	"pecan":    "nutty",
	"cashew":   "nutty",

	// Sweet
	"brown sugar": "sweet",
	"molasses":    "sweet",
	"honey":       "sweet",
	"maple syrup": "sweet",
	"panela":      "sweet",
	"toffee":      "sweet",
	"caramel":     "sweet",
	"butterscotch": "sweet",

	// Floral
	"jasmine":     "floral",
	"rose":        "floral",
	"lavender":    "floral",
	"elderflower": "floral",
	"hibiscus":    "floral",
	"chamomile":   "floral",

	// Spice
	"cinnamon":     "spice",
	"clove":        "spice",
	"cardamom":     "spice",
	"black pepper": "spice",
	"nutmeg":       "spice",
	"ginger":       "spice",
}

// FlavourGroup returns the SCA tier-2 group for a tasting note, or the note
// itself if no group mapping exists.
func FlavourGroup(note string) string {
	if g, ok := flavourGroups[note]; ok {
		return g
	}
	return note
}
