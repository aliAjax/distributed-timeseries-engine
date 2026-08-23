package metric_domain

func PointsInWindow(points []Point, start, end int64) []Point {
	out := make([]Point, 0, len(points))
	for _, point := range points {
		if point.Timestamp >= start && point.Timestamp < end {
			out = append(out, point)
		}
	}
	return out
}
