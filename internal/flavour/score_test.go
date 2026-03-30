package flavour

import (
	"math"
	"testing"
)

func TestComposite_AllPresent(t *testing.T) {
	v := Composite("natural", "light", "ET", "heirloom")
	if v.Fruity < 0.5 {
		t.Errorf("natural+ET+heirloom should be fruity, got %.2f", v.Fruity)
	}
	if v.Floral < 0.5 {
		t.Errorf("ET+heirloom should be floral, got %.2f", v.Floral)
	}
}

func TestComposite_MissingAttributes(t *testing.T) {
	v := Composite("washed", "", "", "")
	if v.Bright < 0.5 {
		t.Errorf("washed alone should retain bright, got %.2f", v.Bright)
	}
}

func TestMatchScore_IdenticalVectors(t *testing.T) {
	v := Vector{Fruity: 0.8, Sweet: 0.3, Bright: 0.7, Body: 0.4, Floral: 0.6, Earthy: 0.1}
	score := MatchScore(v, v)
	if math.Abs(score-1.0) > 0.01 {
		t.Errorf("identical vectors should score 1.0, got %.4f", score)
	}
}

func TestMatchScore_OppositeVectors(t *testing.T) {
	a := Vector{Fruity: 1.0, Sweet: 0.0, Bright: 1.0, Body: 0.0, Floral: 1.0, Earthy: 0.0}
	b := Vector{Fruity: 0.0, Sweet: 1.0, Bright: 0.0, Body: 1.0, Floral: 0.0, Earthy: 1.0}
	score := MatchScore(a, b)
	if score > 0.1 {
		t.Errorf("opposite vectors should score near 0, got %.4f", score)
	}
}

func TestTargetFromAnswers(t *testing.T) {
	target := TargetFromAnswers("fruity", "bright", "light", "floral", "classic")
	if target.Fruity < 0.5 {
		t.Errorf("fruity sweetness should boost fruity, got %.2f", target.Fruity)
	}
	if target.Bright < 0.5 {
		t.Errorf("bright answer should boost bright, got %.2f", target.Bright)
	}
}
