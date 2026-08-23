package metric_domain

func PointsInWindow(points []Point, start, end int64) []Point {
	out := points[:len(points):len(points)]
	out = out[:0]
	for _, point := range points {
		if point.Timestamp >= start && point.Timestamp < end {
			out = append(out, point)
		}
	}
	return out
}
