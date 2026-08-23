package query_parser

import (
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Expr struct {
	Metric   string
	Matchers []metric_domain.LabelMatcher
	Function string
	Range    time.Duration
	GroupBy  []string
}

func Parse(input string) (Expr, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Expr{}, fmt.Errorf("empty query")
	}
	e := Expr{Metric: input}
	if p := strings.Index(input, "["); p >= 0 {
		end := strings.Index(input[p:], "]")
		if end < 0 {
			return e, fmt.Errorf("range selector")
		}
		dur, er := time.ParseDuration(input[p+1 : p+end])
		if er != nil {
			return e, er
		}
		e.Range = dur
		e.Metric = input[:p]
	}
	if p := strings.Index(e.Metric, "{"); p >= 0 {
		end := strings.LastIndex(e.Metric, "}")
		if end < p {
			return e, fmt.Errorf("label selector")
		}
		body := e.Metric[p+1 : end]
		e.Metric = e.Metric[:p]
		for _, part := range strings.Split(body, ",") {
			bits := strings.FieldsFunc(part, func(r rune) bool { return r == '=' || r == '!' || r == '~' })
			if len(bits) < 2 {
				continue
			}
			name, val := strings.TrimSpace(bits[0]), strings.Trim(strings.TrimSpace(bits[len(bits)-1]), "\"")
			op := metric_domain.MatchEqual
			if strings.Contains(part, "!=") {
				op = metric_domain.MatchNotEqual
			} else if strings.Contains(part, "=~") {
				op = metric_domain.MatchRegexp
			} else if strings.Contains(part, "!~") {
				op = metric_domain.MatchNotRegexp
			}
			e.Matchers = append(e.Matchers, metric_domain.LabelMatcher{Name: name, Value: val, Type: op})
		}
	}
	if e.Metric == "" {
		return e, fmt.Errorf("metric required")
	}
	return e, nil
}
func ParseDuration(v string) (time.Duration, error) {
	if n, e := strconv.Atoi(v); e == nil {
		return time.Duration(n) * time.Second, nil
	}
	return time.ParseDuration(v)
}
