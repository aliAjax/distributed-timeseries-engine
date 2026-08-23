package label_index

import (
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"regexp"
	"sort"
	"sync"
)

type Index struct {
	mu     sync.RWMutex
	series map[string]metric_domain.LabelSet
	values map[string]map[string]map[string]struct{}
}

func New() *Index {
	return &Index{series: map[string]metric_domain.LabelSet{}, values: map[string]map[string]map[string]struct{}{}}
}
func (i *Index) Add(id string, l metric_domain.LabelSet) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.series[id] = l.Clone()
	for _, x := range l {
		if i.values[x.Name] == nil {
			i.values[x.Name] = map[string]map[string]struct{}{}
		}
		if i.values[x.Name][x.Value] == nil {
			i.values[x.Name][x.Value] = map[string]struct{}{}
		}
		i.values[x.Name][x.Value][id] = struct{}{}
	}
}
func (i *Index) Match(ms []metric_domain.LabelMatcher) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := make([]string, 0)
	for id, l := range i.series {
		ok := true
		for _, m := range ms {
			if !m.Match(l) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}
func (i *Index) Names() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := make([]string, 0, len(i.values))
	for n := range i.values {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
func (i *Index) Values(name string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var out []string
	for v := range i.values[name] {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
func (i *Index) Cardinality() int { i.mu.RLock(); defer i.mu.RUnlock(); return len(i.series) }
func CompileMatcher(m metric_domain.LabelMatcher) (func(string) bool, error) {
	if m.Type == metric_domain.MatchRegexp || m.Type == metric_domain.MatchNotRegexp {
		r, e := regexp.Compile(m.Value)
		if e != nil {
			return nil, e
		}
		return func(v string) bool {
			ok := r.MatchString(v)
			if m.Type == metric_domain.MatchNotRegexp {
				return !ok
			}
			return ok
		}, nil
	}
	return func(v string) bool {
		if m.Type == metric_domain.MatchNotEqual {
			return v != m.Value
		}
		return v == m.Value
	}, nil
}
