package normalise

import "testing"

func TestNormaliseOrigin(t *testing.T) {
	tests := []struct {
		name        string
		originRaw   string
		regionRaw   string
		wantCountry string
		wantRegion  string
	}{
		{
			name:        "exact country match",
			originRaw:   "Ethiopia",
			wantCountry: "ET",
		},
		{
			name:        "adjective form",
			originRaw:   "Colombian",
			wantCountry: "CO",
		},
		{
			name:        "country with region suffix",
			originRaw:   "Ethiopia Yirgacheffe",
			wantCountry: "ET",
			wantRegion:  "yirgacheffe",
		},
		{
			name:        "country with comma region",
			originRaw:   "Colombia, Huila",
			wantCountry: "CO",
			wantRegion:  "huila",
		},
		{
			name:        "explicit regionRaw overrides extracted",
			originRaw:   "Ethiopia Yirgacheffe",
			regionRaw:   "Gedeo Zone",
			wantCountry: "ET",
			wantRegion:  "Gedeo Zone",
		},
		{
			name:        "region-only origin (well-known region)",
			originRaw:   "Yirgacheffe",
			wantCountry: "ET",
			wantRegion:  "yirgacheffe",
		},
		{
			name:        "PNG abbreviation",
			originRaw:   "PNG",
			wantCountry: "PG",
		},
		{
			name:        "Sumatra maps to Indonesia",
			originRaw:   "Sumatra",
			wantCountry: "ID",
			wantRegion:  "sumatra",
		},
		{
			name:        "case insensitive",
			originRaw:   "KENYA",
			wantCountry: "KE",
		},
		{
			name:        "empty string",
			originRaw:   "",
			wantCountry: "",
			wantRegion:  "",
		},
		{
			name:        "unrecognised origin",
			originRaw:   "Unknown Planet",
			wantCountry: "",
			wantRegion:  "",
		},
		{
			name:        "Costa Rica with Tarrazu",
			originRaw:   "Costa Rica Tarrazu",
			wantCountry: "CR",
			wantRegion:  "tarrazu",
		},
		{
			name:        "region only with explicit regionRaw",
			originRaw:   "Huila",
			regionRaw:   "Huila",
			wantCountry: "CO",
			wantRegion:  "Huila",
		},
		{
			name:        "Yunnan maps to China",
			originRaw:   "Yunnan",
			wantCountry: "CN",
			wantRegion:  "yunnan",
		},
		{
			name:        "Colombia misspelled as Columbia",
			originRaw:   "Columbia",
			wantCountry: "CO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCountry, gotRegion := NormaliseOrigin(tt.originRaw, tt.regionRaw)
			if gotCountry != tt.wantCountry {
				t.Errorf("country: got %q, want %q", gotCountry, tt.wantCountry)
			}
			if gotRegion != tt.wantRegion {
				t.Errorf("region: got %q, want %q", gotRegion, tt.wantRegion)
			}
		})
	}
}

func TestNormaliseProcess(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Washed", ProcessWashed},
		{"natural", ProcessNatural},
		{"HONEY", ProcessHoney},
		{"anaerobic natural", ProcessAnaerobic},
		{"giling basah", ProcessWetHulled},
		{"carbonic maceration", ProcessExperimental},
		{"", ProcessUnknown},
		{"unknown", ProcessUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormaliseProcess(tt.input)
			if got != tt.want {
				t.Errorf("NormaliseProcess(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormaliseRoastLevel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Light", RoastLight},
		{"filter roast", RoastLight},
		{"medium", RoastMedium},
		{"espresso", RoastMedium},
		{"full city", RoastMedDark},
		{"dark", RoastDark},
		{"", RoastUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormaliseRoastLevel(tt.input)
			if got != tt.want {
				t.Errorf("NormaliseRoastLevel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
