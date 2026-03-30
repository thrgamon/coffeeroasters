package flavour

import "testing"

func TestProcessWeights_KnownProcess(t *testing.T) {
	w := ProcessWeights("natural")
	if w.Fruity < 0.7 {
		t.Errorf("natural should have high fruity weight, got %.2f", w.Fruity)
	}
	if w.Body < 0.6 {
		t.Errorf("natural should have high body weight, got %.2f", w.Body)
	}
}

func TestProcessWeights_UnknownProcess(t *testing.T) {
	w := ProcessWeights("unknown-process")
	if w != (Vector{}) {
		t.Errorf("unknown process should return zero vector, got %+v", w)
	}
}

func TestRoastWeights_Light(t *testing.T) {
	w := RoastWeights("light")
	if w.Bright < 0.7 {
		t.Errorf("light roast should have high bright weight, got %.2f", w.Bright)
	}
}

func TestOriginWeights_KnownOrigin(t *testing.T) {
	w := OriginWeights("ET")
	if w.Floral < 0.5 {
		t.Errorf("Ethiopia should have high floral weight, got %.2f", w.Floral)
	}
}

func TestOriginWeights_UnknownOrigin(t *testing.T) {
	w := OriginWeights("XX")
	if w != (Vector{}) {
		t.Errorf("unknown origin should return zero vector, got %+v", w)
	}
}
