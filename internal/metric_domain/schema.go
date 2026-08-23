package metric_domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type LabelSchema struct {
	Name           string
	Required       bool
	Pattern        string
	MaxValueLength int
}
type MetricSchema struct {
	Tenant, Metric, Unit string
	Labels               []LabelSchema
	UpdatedAt            time.Time
	Version              int
}

func (s MetricSchema) Validate() error {
	if s.Tenant == "" || s.Metric == "" {
		return fmt.Errorf("schema tenant and metric required")
	}
	for _, l := range s.Labels {
		if l.Name == "" {
			return fmt.Errorf("label name required")
		}
		if l.Pattern != "" {
			if _, e := regexp.Compile(l.Pattern); e != nil {
				return e
			}
		}
	}
	return nil
}
func (s MetricSchema) Check(l LabelSet) error {
	if e := s.Validate(); e != nil {
		return e
	}
	for _, rule := range s.Labels {
		v, ok := l.Get(rule.Name)
		if rule.Required && !ok {
			return fmt.Errorf("required label %s missing", rule.Name)
		}
		if ok && rule.MaxValueLength > 0 && len(v) > rule.MaxValueLength {
			return fmt.Errorf("label %s too long", rule.Name)
		}
		if ok && rule.Pattern != "" {
			ok2, _ := regexp.MatchString(rule.Pattern, v)
			if !ok2 {
				return fmt.Errorf("label %s pattern mismatch", rule.Name)
			}
		}
	}
	return nil
}
func NormalizeUnit(unit string) string { return strings.ToLower(strings.TrimSpace(unit)) }
