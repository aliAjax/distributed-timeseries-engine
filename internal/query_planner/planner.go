package query_planner

import (
	"errors"
	"example.com/distributed-timeseries-engine/internal/query_parser"
	"fmt"
	"time"
)

var ErrInvalidPlan = errors.New("invalid query plan")

type PlanFailure struct {
	Start  time.Time
	End    time.Time
	Detail string
}

func (e PlanFailure) Error() string {
	return fmt.Sprintf("plan range %s..%s: %s", e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339), e.Detail)
}
func (e PlanFailure) Unwrap() error { return nil }

func invalidPlan(start, end time.Time, err error) error {
	return PlanFailure{Start: start, End: end, Detail: err.Error()}
}

type Plan struct {
	Expr                    query_parser.Expr
	Start, End              time.Time
	SeriesLimit, PointLimit int
	Parallelism             int
}

func Build(e query_parser.Expr, start, end time.Time, limits ...int) (Plan, error) {
	if !end.After(start) {
		return Plan{}, invalidPlan(start, end, errors.New("end must follow start"))
	}
	p := Plan{Expr: e, Start: start, End: end, SeriesLimit: 10000, PointLimit: 100000, Parallelism: 4}
	if len(limits) > 0 && limits[0] > 0 {
		p.PointLimit = limits[0]
	}
	return p, nil
}
func Cost(p Plan, series, points int) float64 {
	if series == 0 {
		return 0
	}
	return float64(points) + float64(series)*10 + float64(p.End.Sub(p.Start))/float64(time.Minute)
}
