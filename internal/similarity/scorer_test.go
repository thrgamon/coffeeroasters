package similarity

import (
	"math"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func TestScore_IdenticalCoffees(t *testing.T) {
	a := CoffeeAttrs{
		TastingNotes: []string{"blueberry", "chocolate", "caramel"},
		Process:      "natural",
		RoastLevel:   "light",
		Variety:      "bourbon",
		RegionID:     ptr(int32(1)),
	}
	score := Score(a, a)
	if score < 0.99 {
		t.Errorf("identical coffees should score ~1.0, got %.4f", score)
	}
}

func TestScore_CompletelyDifferent(t *testing.T) {
	a := CoffeeAttrs{
		TastingNotes: []string{"blueberry", "jasmine"},
		Process:      "natural",
		RoastLevel:   "light",
		Variety:      "gesha",
		RegionID:     ptr(int32(1)),
	}
	b := CoffeeAttrs{
		TastingNotes: []string{"dark chocolate", "walnut"},
		Process:      "washed",
		RoastLevel:   "dark",
		Variety:      "catimor",
		RegionID:     ptr(int32(99)),
	}
	score := Score(a, b)
	if score > 0.15 {
		t.Errorf("very different coffees should score low, got %.4f", score)
	}
}

func TestScore_FlavourGroupPartialCredit(t *testing.T) {
	a := CoffeeAttrs{
		TastingNotes: []string{"blueberry", "chocolate"},
	}
	b := CoffeeAttrs{
		TastingNotes: []string{"strawberry", "chocolate"},
	}
	score := Score(a, b)

	// "chocolate" is exact match. "blueberry" and "strawberry" share "berry" group.
	// Union = 3 (blueberry, strawberry, chocolate)
	// Exact = 1 (chocolate), Group = 1 (blueberry->berry matches strawberry->berry)
	// Note similarity = (1 + 0.5) / 3 = 0.5
	// Total = 0.5 * 0.40 = 0.20
	if score < 0.15 || score > 0.25 {
		t.Errorf("partial flavour group credit score = %.4f, expected ~0.20", score)
	}
}

func TestScore_RoastOneStepApart(t *testing.T) {
	a := CoffeeAttrs{RoastLevel: "light"}
	b := CoffeeAttrs{RoastLevel: "medium-light"}
	score := Score(a, b)

	// Only roast contributes: 0.5 * 0.15 = 0.075
	expected := 0.5 * weightRoast
	if math.Abs(score-expected) > 0.001 {
		t.Errorf("one-step roast score = %.4f, expected %.4f", score, expected)
	}
}

func TestScore_RoastTwoStepsApart(t *testing.T) {
	a := CoffeeAttrs{RoastLevel: "light"}
	b := CoffeeAttrs{RoastLevel: "medium"}
	score := Score(a, b)
	if score > 0.001 {
		t.Errorf("two-step roast should contribute 0, got %.4f", score)
	}
}

func TestRank_ExcludesSourceAndLowScores(t *testing.T) {
	source := CoffeeAttrs{
		CoffeeID:     1,
		TastingNotes: []string{"blueberry", "chocolate"},
		Process:      "natural",
		RoastLevel:   "light",
		Variety:      "bourbon",
		RegionID:     ptr(int32(1)),
	}

	candidates := []CoffeeAttrs{
		source, // should be excluded (same ID)
		{
			CoffeeID:     2,
			TastingNotes: []string{"blueberry", "chocolate"},
			Process:      "natural",
			RoastLevel:   "light",
			Variety:      "bourbon",
			RegionID:     ptr(int32(1)),
		},
		{
			CoffeeID:     3,
			TastingNotes: []string{"walnut"},
			Process:      "washed",
			RoastLevel:   "dark",
			Variety:      "catimor",
			RegionID:     ptr(int32(99)),
		},
	}

	results := Rank(source, candidates, 10)

	for _, r := range results {
		if r.CoffeeID == 1 {
			t.Error("source coffee should be excluded from results")
		}
	}

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}

	if results[0].CoffeeID != 2 {
		t.Errorf("expected coffee 2 to rank first, got %d", results[0].CoffeeID)
	}
}

func TestRank_RespectsLimit(t *testing.T) {
	source := CoffeeAttrs{
		CoffeeID:     1,
		TastingNotes: []string{"chocolate"},
		Process:      "natural",
	}

	var candidates []CoffeeAttrs
	for i := int64(2); i <= 20; i++ {
		candidates = append(candidates, CoffeeAttrs{
			CoffeeID:     i,
			TastingNotes: []string{"chocolate"},
			Process:      "natural",
		})
	}

	results := Rank(source, candidates, 5)
	if len(results) > 5 {
		t.Errorf("expected max 5 results, got %d", len(results))
	}
}

func TestHaversineKm(t *testing.T) {
	// London to Paris ~344km
	dist := haversineKm(51.5074, -0.1278, 48.8566, 2.3522)
	if dist < 340 || dist > 350 {
		t.Errorf("London-Paris distance = %.1f km, expected ~344", dist)
	}
}
