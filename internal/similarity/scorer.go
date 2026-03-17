package similarity

import (
	"math"
	"sort"
)

// Weights for each similarity dimension.
const (
	weightTastingNotes = 0.40
	weightProcess      = 0.15
	weightRoast        = 0.15
	weightVariety      = 0.15
	weightRegion       = 0.15

	// Minimum score threshold for inclusion in results.
	minScoreThreshold = 0.15

	// Maximum distance in km for region proximity scoring.
	maxRegionDistanceKm = 500.0
)

// ScoredCoffee pairs a coffee ID with its similarity score.
type ScoredCoffee struct {
	CoffeeID int64
	Score    float64
}

// CoffeeAttrs holds the attributes used for similarity scoring.
type CoffeeAttrs struct {
	CoffeeID     int64
	TastingNotes []string
	Process      string
	RoastLevel   string
	Variety      string
	RegionID     *int32
	Latitude     *float64
	Longitude    *float64
}

// roastOrdinal maps canonical roast levels to an ordinal value for distance
// calculation.
var roastOrdinal = map[string]int{
	"light":       0,
	"medium-light": 1,
	"medium":      2,
	"medium-dark": 3,
	"dark":        4,
}

// Score computes a weighted similarity between a source coffee and a candidate.
func Score(source, candidate CoffeeAttrs) float64 {
	var total float64

	total += weightTastingNotes * tastingNoteSimilarity(source.TastingNotes, candidate.TastingNotes)
	total += weightProcess * exactMatch(source.Process, candidate.Process)
	total += weightRoast * roastSimilarity(source.RoastLevel, candidate.RoastLevel)
	total += weightVariety * varietySimilarity(source.Variety, candidate.Variety)
	total += weightRegion * regionSimilarity(source, candidate)

	return total
}

// Rank scores all candidates against the source and returns the top N sorted
// by score descending. Excludes scores below the minimum threshold and the
// source coffee itself.
func Rank(source CoffeeAttrs, candidates []CoffeeAttrs, limit int) []ScoredCoffee {
	var scored []ScoredCoffee

	for _, c := range candidates {
		if c.CoffeeID == source.CoffeeID {
			continue
		}
		s := Score(source, c)
		if s >= minScoreThreshold {
			scored = append(scored, ScoredCoffee{CoffeeID: c.CoffeeID, Score: s})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	return scored
}

func tastingNoteSimilarity(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	bSet := make(map[string]bool, len(b))
	bGroups := make(map[string]bool, len(b))
	for _, n := range b {
		bSet[n] = true
		bGroups[FlavourGroup(n)] = true
	}

	var exactMatches, groupMatches float64
	aGroups := make(map[string]bool, len(a))
	for _, n := range a {
		aGroups[FlavourGroup(n)] = true
		if bSet[n] {
			exactMatches++
		}
	}

	// Count group-level matches (excluding notes that already matched exactly)
	for _, n := range a {
		if bSet[n] {
			continue
		}
		g := FlavourGroup(n)
		if bGroups[g] && g != n {
			groupMatches++
		}
	}

	// Union size for Jaccard-like calculation
	union := float64(countUnion(a, b))
	if union == 0 {
		return 0
	}

	// Exact matches count fully, group matches count at 0.5
	return (exactMatches + groupMatches*0.5) / union
}

func countUnion(a, b []string) int {
	seen := make(map[string]bool, len(a)+len(b))
	for _, n := range a {
		seen[n] = true
	}
	for _, n := range b {
		seen[n] = true
	}
	return len(seen)
}

func exactMatch(a, b string) float64 {
	if a != "" && a == b {
		return 1.0
	}
	return 0.0
}

func roastSimilarity(a, b string) float64 {
	if a == "" || b == "" {
		return 0
	}
	ordA, okA := roastOrdinal[a]
	ordB, okB := roastOrdinal[b]
	if !okA || !okB {
		return 0
	}
	diff := abs(ordA - ordB)
	switch diff {
	case 0:
		return 1.0
	case 1:
		return 0.5
	default:
		return 0.0
	}
}

func varietySimilarity(a, b string) float64 {
	if a != "" && a == b {
		return 1.0
	}
	return 0.0
}

func regionSimilarity(a, b CoffeeAttrs) float64 {
	// Same region ID is a perfect match
	if a.RegionID != nil && b.RegionID != nil && *a.RegionID == *b.RegionID {
		return 1.0
	}

	// Fall back to Haversine distance if coordinates available
	if a.Latitude == nil || a.Longitude == nil || b.Latitude == nil || b.Longitude == nil {
		return 0.0
	}

	dist := haversineKm(*a.Latitude, *a.Longitude, *b.Latitude, *b.Longitude)
	if dist >= maxRegionDistanceKm {
		return 0.0
	}

	return 1.0 - dist/maxRegionDistanceKm
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
