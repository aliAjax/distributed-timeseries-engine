package metric_domain

func SelectQuality(points []Point, minimum Quality) []Point {
	out := points[:len(points):len(points)]
	out = out[:0]
	for _, point := range points {
		if point.Quality >= minimum {
			out = append(out, point)
		}
	}
	return out
}
