package metric_domain

func QualityScore(q Quality) float64 {
	switch q {
	case QualityGood:
		return 1
	case QualityStale:
		return .5
	case QualityBad:
		return 0
	default:
		return .25
	}
}
func MergeQuality(values []Quality) Quality {
	best := QualityGood
	for _, q := range values {
		if QualityScore(q) < QualityScore(best) {
			best = q
		}
	}
	return best
}
func FilterQuality(points []Point, min Quality) []Point {
	out := make([]Point, 0, len(points))
	for _, p := range points {
		if p.Quality >= min {
			out = append(out, p)
		}
	}
	return out
}
