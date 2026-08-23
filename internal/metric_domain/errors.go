package metric_domain

import "fmt"

type IngestError struct {
	Index int
	Err   error
}

func (e IngestError) Error() string { return fmt.Sprintf("sample %d: %v", e.Index, e.Err) }
func (e IngestError) Unwrap() error { return e.Err }

type Point struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	Quality   Quality `json:"quality"`
}
type Series struct {
	Tenant, Metric string
	Labels         LabelSet
	Points         []Point
}

func (s Series) Sort() {
	for i := 1; i < len(s.Points); i++ {
		x := s.Points[i]
		j := i - 1
		for j >= 0 && s.Points[j].Timestamp > x.Timestamp {
			s.Points[j+1] = s.Points[j]
			j--
		}
		s.Points[j+1] = x
	}
}
func (s Series) Deduplicate() Series {
	s.Sort()
	out := Series{Tenant: s.Tenant, Metric: s.Metric, Labels: s.Labels, Points: make([]Point, 0, len(s.Points))}
	for _, p := range s.Points {
		if len(out.Points) > 0 && out.Points[len(out.Points)-1].Timestamp == p.Timestamp {
			if p.Quality >= out.Points[len(out.Points)-1].Quality {
				out.Points[len(out.Points)-1] = p
			}
			continue
		}
		out.Points = append(out.Points, p)
	}
	return out
}
