package normalise

import "testing"

func TestNormaliseVariety(t *testing.T) {
	tests := []struct {
		raw         string
		wantVariety string
		wantSpecies string
	}{
		// Exact matches
		{"Bourbon", "bourbon", "arabica"},
		{"typica", "typica", "arabica"},
		{"SL28", "sl28", "arabica"},
		{"SL-28", "sl28", "arabica"},

		// Alias matches
		{"Geisha", "gesha", "arabica"},
		{"gesha", "gesha", "arabica"},
		{"Heirloom Varieties", "heirloom", "arabica"},
		{"Ethiopian Heirloom", "heirloom", "arabica"},
		{"Bourbon Rouge", "red-bourbon", "arabica"},
		{"Caturra Rojo", "caturra", "arabica"},
		{"Yellow Bourbon", "yellow-bourbon", "arabica"},

		// Multi-variety
		{"Bourbon, Caturra", "bourbon,caturra", "arabica"},
		{"SL28 / SL34", "sl28,sl34", "arabica"},
		{"Typica & Bourbon", "typica,bourbon", "arabica"},

		// Unknown -> empty for LLM classifier
		{"some weird variety", "", ""},
		{"", "", ""},

		// Ethiopian selections
		{"74110", "74110", "arabica"},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			gotVariety, gotSpecies := NormaliseVariety(tt.raw)
			if gotVariety != tt.wantVariety {
				t.Errorf("NormaliseVariety(%q) variety = %q, want %q", tt.raw, gotVariety, tt.wantVariety)
			}
			if gotSpecies != tt.wantSpecies {
				t.Errorf("NormaliseVariety(%q) species = %q, want %q", tt.raw, gotSpecies, tt.wantSpecies)
			}
		})
	}
}
