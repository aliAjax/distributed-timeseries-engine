package query_executor

import (
	"context"
	"example.com/distributed-timeseries-engine/internal/aggregation"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/query_planner"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Reader interface {
	QuerySeries(context.Context, string, []metric_domain.LabelMatcher, time.Time, time.Time) (map[string][]metric_domain.Point, error)
}
type ResultSeries struct {
	ID         string                `json:"id"`
	Points     []metric_domain.Point `json:"points"`
	Aggregates []AggregatePoint      `json:"aggregates,omitempty"`
}
type AggregatePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	Count     int     `json:"count"`
}
type Result struct {
	Series    []ResultSeries `json:"series"`
	Examined  int            `json:"examined"`
	Truncated bool           `json:"truncated"`
}
type Executor struct {
	Reader     Reader
	MaxWorkers int
	MaxPoints  int
}

func (e Executor) Execute(ctx context.Context, p query_planner.Plan) (Result, error) {
	if e.Reader == nil {
		return Result{}, fmt.Errorf("reader is required")
	}
	if e.MaxWorkers < 1 {
		e.MaxWorkers = 4
	}
	if e.MaxPoints < 1 {
		e.MaxPoints = 100000
	}
	data, err := e.Reader.QuerySeries(ctx, p.Expr.Metric, p.Expr.Matchers, p.Start, p.End)
	if err != nil {
		return Result{}, err
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := Result{Series: make([]ResultSeries, 0, len(keys))}
	jobs := make(chan string)
	out := make(chan ResultSeries, len(keys))
	var wg sync.WaitGroup
	workers := e.MaxWorkers
	if workers > len(keys) {
		workers = len(keys)
	}
	if workers == 0 {
		return res, nil
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				pts := data[id]
				rs := ResultSeries{ID: id, Points: append([]metric_domain.Point(nil), pts...)}
				if p.Expr.Function != "" {
					step := time.Minute
					if p.Expr.Range > 0 && p.Expr.Range < time.Hour {
						step = p.Expr.Range / 10
						if step < time.Second {
							step = time.Second
						}
					}
					for _, w := range aggregation.Windows(pts, p.Start, p.End, step) {
						rs.Aggregates = append(rs.Aggregates, AggregatePoint{Timestamp: w.Start.UnixNano(), Value: aggregation.Value(w.Points, p.Expr.Function), Count: len(w.Points)})
					}
				}
				out <- rs
			}
		}()
	}
	go func() {
		for _, k := range keys {
			jobs <- k
		}
		close(jobs)
		wg.Wait()
		close(out)
	}()
	for x := range out {
		res.Examined += len(x.Points)
		if res.Examined > e.MaxPoints {
			res.Truncated = true
			remain := e.MaxPoints - (res.Examined - len(x.Points))
			if remain < 0 {
				remain = 0
			}
			x.Points = x.Points[:min(remain, len(x.Points))]
			res.Examined = e.MaxPoints
		}
		res.Series = append(res.Series, x)
		if res.Truncated {
			break
		}
	}
	return res, nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (e Executor) Aggregate(points []metric_domain.Point, op string, start, end time.Time, step time.Duration) []AggregatePoint {
	var out []AggregatePoint
	for _, w := range aggregation.Windows(points, start, end, step) {
		out = append(out, AggregatePoint{Timestamp: w.Start.UnixNano(), Value: aggregation.Value(w.Points, op), Count: len(w.Points)})
	}
	return out
}
func Merge(results ...Result) Result {
	merged := Result{}
	byID := map[string]ResultSeries{}
	for _, r := range results {
		merged.Examined += r.Examined
		merged.Truncated = merged.Truncated || r.Truncated
		for _, s := range r.Series {
			x := byID[s.ID]
			x.ID = s.ID
			x.Points = append(x.Points, s.Points...)
			x.Aggregates = append(x.Aggregates, s.Aggregates...)
			byID[s.ID] = x
		}
	}
	for _, s := range byID {
		sort.Slice(s.Points, func(i, j int) bool { return s.Points[i].Timestamp < s.Points[j].Timestamp })
		merged.Series = append(merged.Series, s)
	}
	sort.Slice(merged.Series, func(i, j int) bool { return merged.Series[i].ID < merged.Series[j].ID })
	return merged
}
