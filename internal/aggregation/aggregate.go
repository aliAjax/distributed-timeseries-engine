package aggregation

import (
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"math"
	"sort"
	"time"
)

type Window struct {
	Start, End time.Time
	Points     []metric_domain.Point
}

func Windows(points []metric_domain.Point, start, end time.Time, step time.Duration) []Window {
	if step <= 0 {
		return nil
	}
	var out []Window
	for t := start; t.Before(end); t = t.Add(step) {
		w := Window{Start: t, End: t.Add(step)}
		for _, p := range points {
			pt := time.Unix(0, p.Timestamp)
			if !pt.Before(w.Start) && pt.Before(w.End) {
				w.Points = append(w.Points, p)
			}
		}
		out = append(out, w)
	}
	return out
}
func Value(points []metric_domain.Point, op string) float64 {
	if len(points) == 0 {
		return math.NaN()
	}
	switch op {
	case "sum":
		var n float64
		for _, p := range points {
			n += p.Value
		}
		return n
	case "min":
		n := points[0].Value
		for _, p := range points[1:] {
			if p.Value < n {
				n = p.Value
			}
		}
		return n
	case "max":
		n := points[0].Value
		for _, p := range points[1:] {
			if p.Value > n {
				n = p.Value
			}
		}
		return n
	case "count":
		return float64(len(points))
	default:
		var n float64
		for _, p := range points {
			n += p.Value
		}
		return n / float64(len(points))
	}
}
func Rate(points []metric_domain.Point) float64 {
	if len(points) < 2 {
		return 0
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Timestamp < points[j].Timestamp })
	dt := float64(points[len(points)-1].Timestamp-points[0].Timestamp) / 1e9
	if dt <= 0 {
		return 0
	}
	return (points[len(points)-1].Value - points[0].Value) / dt
}
func Quantile(points []metric_domain.Point, q float64) float64 {
	if len(points) == 0 {
		return math.NaN()
	}
	if q < 0 {
		q = 0
	}
	if q > 1 {
		q = 1
	}
	v := make([]float64, len(points))
	for j, p := range points {
		v[j] = p.Value
	}
	sort.Float64s(v)
	idx := int(math.Round(q * float64(len(v)-1)))
	return v[idx]
}
