package metric_domain

import "sort"

type Window struct {
	Start, End int64
	Points     []Point
}

func NewWindow(start, end int64) Window {
	return Window{Start: start, End: end, Points: make([]Point, 0)}
}
func (w *Window) Add(point Point) bool {
	if point.Timestamp < w.Start || point.Timestamp >= w.End {
		return false
	}
	w.Points = append(w.Points, point)
	return true
}
func (w *Window) Sort() {
	sort.SliceStable(w.Points, func(i, j int) bool { return w.Points[i].Timestamp < w.Points[j].Timestamp })
}
func (w Window) Values() []float64 {
	out := make([]float64, len(w.Points))
	for i, p := range w.Points {
		out[i] = p.Value
	}
	return out
}
func (w Window) Mean() float64 {
	if len(w.Points) == 0 {
		return 0
	}
	var total float64
	for _, p := range w.Points {
		total += p.Value
	}
	return total / float64(len(w.Points))
}
func (w Window) MinMax() (float64, float64) {
	if len(w.Points) == 0 {
		return 0, 0
	}
	min, max := w.Points[0].Value, w.Points[0].Value
	for _, p := range w.Points[1:] {
		if p.Value < min {
			min = p.Value
		}
		if p.Value > max {
			max = p.Value
		}
	}
	return min, max
}
