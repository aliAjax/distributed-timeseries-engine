package metric_domain

import "sort"

type Bucket struct {
	Start int64
	End   int64
	Count int
	Sum   float64
	Min   float64
	Max   float64
}

func Resample(points []Point, start, end, step int64) []Bucket {
	if step <= 0 {
		return nil
	}
	count := int((end - start + step - 1) / step)
	out := make([]Bucket, count)
	for i := range out {
		out[i] = Bucket{Start: start + int64(i)*step, End: start + int64(i+1)*step, Min: 0, Max: 0}
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Timestamp < points[j].Timestamp })
	for _, p := range points {
		if p.Timestamp < start || p.Timestamp >= end {
			continue
		}
		i := int((p.Timestamp - start) / step)
		b := &out[i]
		if b.Count == 0 {
			b.Min, b.Max = p.Value, p.Value
		} else {
			if p.Value < b.Min {
				b.Min = p.Value
			}
			if p.Value > b.Max {
				b.Max = p.Value
			}
		}
		b.Count++
		b.Sum += p.Value
	}
	return out
}
func (b Bucket) Mean() float64 {
	if b.Count == 0 {
		return 0
	}
	return b.Sum / float64(b.Count)
}
