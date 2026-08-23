package metric_domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	ErrInvalidMetric = errors.New("invalid metric")
	ErrInvalidSample = errors.New("invalid sample")
	ErrOutOfWindow   = errors.New("sample outside accepted window")
)

type Quality uint8

const (
	QualityUnknown Quality = iota
	QualityGood
	QualityStale
	QualityBad
)

func (q Quality) String() string {
	switch q {
	case QualityGood:
		return "good"
	case QualityStale:
		return "stale"
	case QualityBad:
		return "bad"
	default:
		return "unknown"
	}
}

type Label struct{ Name, Value string }
type LabelSet []Label

func (l LabelSet) Clone() LabelSet { out := make(LabelSet, len(l)); copy(out, l); return out }
func (l LabelSet) Canonical() string {
	c := l.Clone()
	sort.Slice(c, func(i, j int) bool { return c[i].Name < c[j].Name })
	var b strings.Builder
	for _, x := range c {
		b.WriteString(x.Name)
		b.WriteByte('=')
		b.WriteString(x.Value)
		b.WriteByte(',')
	}
	return b.String()
}
func (l LabelSet) Get(name string) (string, bool) {
	for _, x := range l {
		if x.Name == name {
			return x.Value, true
		}
	}
	return "", false
}
func (l LabelSet) Valid() bool {
	seen := map[string]bool{}
	for _, x := range l {
		if x.Name == "" || seen[x.Name] {
			return false
		}
		seen[x.Name] = true
	}
	return true
}

type Metric struct {
	Tenant, Name, Unit string
	Labels             LabelSet
	CreatedAt          time.Time
	Description        string
}

func (m Metric) SeriesID() string { return m.Tenant + "/" + m.Name + "{" + m.Labels.Canonical() + "}" }
func (m Metric) Validate() error {
	if m.Tenant == "" || m.Name == "" || !m.Labels.Valid() {
		return ErrInvalidMetric
	}
	return nil
}

type Sample struct {
	Tenant, Metric string
	Labels         LabelSet
	Timestamp      time.Time
	Value          float64
	Quality        Quality
	Unit           string
}

func (s Sample) Validate() error {
	if s.Tenant == "" || s.Metric == "" || s.Timestamp.IsZero() || !s.Labels.Valid() {
		return ErrInvalidSample
	}
	return nil
}
func (s Sample) SeriesID() string {
	return s.Tenant + "/" + s.Metric + "{" + s.Labels.Canonical() + "}"
}

type TimeWindow struct{ Start, End time.Time }

func (w TimeWindow) Valid() bool               { return !w.Start.IsZero() && !w.End.IsZero() && w.End.After(w.Start) }
func (w TimeWindow) Contains(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }
func (w TimeWindow) Duration() time.Duration   { return w.End.Sub(w.Start) }

type RetentionPolicy struct {
	Name            string
	Raw             time.Duration
	DownsampleAfter time.Duration
	Resolution      time.Duration
	DeleteDelay     time.Duration
}

func (p RetentionPolicy) Validate() error {
	if p.Name == "" || p.Raw <= 0 || p.DeleteDelay < 0 {
		return fmt.Errorf("invalid retention policy")
	}
	if p.Resolution < 0 {
		return fmt.Errorf("invalid resolution")
	}
	return nil
}

type DownsamplingPolicy struct {
	Name       string
	Resolution time.Duration
	Aggregator string
	Enabled    bool
}

func (p DownsamplingPolicy) Validate() error {
	if p.Name == "" || p.Resolution <= 0 {
		return fmt.Errorf("invalid downsampling policy")
	}
	switch p.Aggregator {
	case "avg", "sum", "min", "max", "count":
	default:
		return fmt.Errorf("unsupported aggregator")
	}
	return nil
}

type MatcherType string

const (
	MatchEqual     MatcherType = "="
	MatchNotEqual  MatcherType = "!="
	MatchRegexp    MatcherType = "=~"
	MatchNotRegexp MatcherType = "!~"
)

type LabelMatcher struct {
	Name, Value string
	Type        MatcherType
}

func (m LabelMatcher) Match(labels LabelSet) bool {
	v, ok := labels.Get(m.Name)
	switch m.Type {
	case MatchEqual:
		return ok && v == m.Value
	case MatchNotEqual:
		return !ok || v != m.Value
	case MatchRegexp:
		return ok && strings.Contains(v, m.Value)
	case MatchNotRegexp:
		return !ok || !strings.Contains(v, m.Value)
	default:
		return false
	}
}
