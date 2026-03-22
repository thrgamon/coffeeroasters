package embedding

import (
	"math"
	"testing"
)

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	a := []float64{1, 2, 3}
	sim := CosineSimilarity(a, a)
	if math.Abs(sim-1.0) > 0.0001 {
		t.Errorf("identical vectors: got %.4f, want 1.0", sim)
	}
}

func TestCosineSimilarity_Orthogonal(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{0, 1, 0}
	sim := CosineSimilarity(a, b)
	if math.Abs(sim) > 0.0001 {
		t.Errorf("orthogonal vectors: got %.4f, want 0.0", sim)
	}
}

func TestCosineSimilarity_Opposite(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{-1, -2, -3}
	sim := CosineSimilarity(a, b)
	if math.Abs(sim-(-1.0)) > 0.0001 {
		t.Errorf("opposite vectors: got %.4f, want -1.0", sim)
	}
}

func TestCosineSimilarity_Empty(t *testing.T) {
	if sim := CosineSimilarity(nil, []float64{1}); sim != 0 {
		t.Errorf("nil a: got %.4f, want 0.0", sim)
	}
	if sim := CosineSimilarity([]float64{1}, nil); sim != 0 {
		t.Errorf("nil b: got %.4f, want 0.0", sim)
	}
	if sim := CosineSimilarity([]float64{}, []float64{}); sim != 0 {
		t.Errorf("empty: got %.4f, want 0.0", sim)
	}
}

func TestCosineSimilarity_MismatchedLengths(t *testing.T) {
	a := []float64{1, 2}
	b := []float64{1, 2, 3}
	if sim := CosineSimilarity(a, b); sim != 0 {
		t.Errorf("mismatched lengths: got %.4f, want 0.0", sim)
	}
}
