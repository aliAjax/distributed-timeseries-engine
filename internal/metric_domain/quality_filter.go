package metric_domain

func SelectQuality(points []Point, minimum Quality) []Point {
	out := make([]Point, 0, len(points))
	for _, point := range points {
		if point.Quality >= minimum {
			out = append(out, point)
		}
	}
	return out
}
