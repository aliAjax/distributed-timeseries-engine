package query_planner

import (
	"example.com/distributed-timeseries-engine/internal/query_parser"
	"fmt"
	"time"
)

type Plan struct {
	Expr                    query_parser.Expr
	Start, End              time.Time
	SeriesLimit, PointLimit int
	Parallelism             int
}

func Build(e query_parser.Expr, start, end time.Time, limits ...int) (Plan, error) {
	if !end.After(start) {
		return Plan{}, fmt.Errorf("invalid time range")
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
