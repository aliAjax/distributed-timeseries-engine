package chunk_codec

import (
	"math"
	"sort"
)

type Stats struct {
	Count            int
	Min, Max, Sum    float64
	CompressionRatio float64
}

func Analyze(points []Point, encodedBytes int) Stats {
	st := Stats{Count: len(points), Min: math.Inf(1), Max: math.Inf(-1)}
	for _, p := range points {
		if p.Value < st.Min {
			st.Min = p.Value
		}
		if p.Value > st.Max {
			st.Max = p.Value
		}
		st.Sum += p.Value
	}
	if len(points) == 0 {
		st.Min = 0
		st.Max = 0
	}
	raw := len(points) * 17
	if encodedBytes > 0 {
		st.CompressionRatio = float64(raw) / float64(encodedBytes)
	}
	return st
}
func Downsample(points []Point, step int64) []Point {
	if step <= 0 {
		return append([]Point(nil), points...)
	}
	sorted := make([]Point, len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Timestamp < sorted[j].Timestamp })
	var out []Point
	var bucket int64
	var sum float64
	var count int
	for _, p := range sorted {
		b := p.Timestamp / step
		if count > 0 && b != bucket {
			out = append(out, Point{Timestamp: bucket * step, Value: sum / float64(count), Quality: 1})
			sum = 0
			count = 0
		}
		bucket = b
		sum += p.Value
		count++
	}
	if count > 0 {
		out = append(out, Point{Timestamp: bucket * step, Value: sum / float64(count), Quality: 1})
	}
	return out
}
