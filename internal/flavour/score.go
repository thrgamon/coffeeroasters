package flavour

import "math"

// Composite computes the average flavour vector for a coffee's attributes.
// Empty strings are skipped.
func Composite(process, roast, origin, variety string) Vector {
	var vectors []Vector
	if v := ProcessWeights(process); v != (Vector{}) {
		vectors = append(vectors, v)
	}
	if v := RoastWeights(roast); v != (Vector{}) {
		vectors = append(vectors, v)
	}
	if v := OriginWeights(origin); v != (Vector{}) {
		vectors = append(vectors, v)
	}
	if v := VarietyWeights(variety); v != (Vector{}) {
		vectors = append(vectors, v)
	}
	if len(vectors) == 0 {
		return Vector{}
	}
	var sum Vector
	for _, v := range vectors {
		sum.Fruity += v.Fruity
		sum.Sweet += v.Sweet
		sum.Bright += v.Bright
		sum.Body += v.Body
		sum.Floral += v.Floral
		sum.Earthy += v.Earthy
	}
	n := float64(len(vectors))
	return Vector{
		Fruity: sum.Fruity / n,
		Sweet:  sum.Sweet / n,
		Bright: sum.Bright / n,
		Body:   sum.Body / n,
		Floral: sum.Floral / n,
		Earthy: sum.Earthy / n,
	}
}

// MatchScore returns the cosine similarity between two flavour vectors (0 to 1).
func MatchScore(target, candidate Vector) float64 {
	dot := target.Fruity*candidate.Fruity +
		target.Sweet*candidate.Sweet +
		target.Bright*candidate.Bright +
		target.Body*candidate.Body +
		target.Floral*candidate.Floral +
		target.Earthy*candidate.Earthy

	magA := math.Sqrt(target.Fruity*target.Fruity +
		target.Sweet*target.Sweet +
		target.Bright*target.Bright +
		target.Body*target.Body +
		target.Floral*target.Floral +
		target.Earthy*target.Earthy)

	magB := math.Sqrt(candidate.Fruity*candidate.Fruity +
		candidate.Sweet*candidate.Sweet +
		candidate.Bright*candidate.Bright +
		candidate.Body*candidate.Body +
		candidate.Floral*candidate.Floral +
		candidate.Earthy*candidate.Earthy)

	if magA == 0 || magB == 0 {
		return 0
	}
	return dot / (magA * magB)
}

// TargetFromAnswers builds a target flavour vector from questionnaire answers.
func TargetFromAnswers(sweetness, brightness, body, appeal, adventurous string) Vector {
	var v Vector

	switch sweetness {
	case "fruity":
		v.Fruity = 0.8
		v.Sweet = 0.3
	case "caramel":
		v.Fruity = 0.2
		v.Sweet = 0.8
	default: // "both"
		v.Fruity = 0.5
		v.Sweet = 0.5
	}

	switch brightness {
	case "bright":
		v.Bright = 0.8
	case "smooth":
		v.Bright = 0.2
		v.Body += 0.2
	default:
		v.Bright = 0.5
	}

	switch body {
	case "light":
		v.Body += 0.2
	case "full":
		v.Body += 0.8
	default:
		v.Body += 0.5
	}

	switch appeal {
	case "floral":
		v.Floral = 0.9
	case "chocolate":
		v.Sweet += 0.3
		v.Earthy = 0.2
	case "berry":
		v.Fruity += 0.3
	case "earthy":
		v.Earthy = 0.8
		v.Body += 0.2
	}

	v.Fruity = clamp(v.Fruity)
	v.Sweet = clamp(v.Sweet)
	v.Bright = clamp(v.Bright)
	v.Body = clamp(v.Body)
	v.Floral = clamp(v.Floral)
	v.Earthy = clamp(v.Earthy)

	return v
}

// IsExperimental returns true if a process is anaerobic or experimental.
func IsExperimental(process string) bool {
	return process == "anaerobic" || process == "experimental"
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
