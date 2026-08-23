package repository

import (
	"context"
	"example.com/distributed-timeseries-engine/internal/block_store"
	"example.com/distributed-timeseries-engine/internal/chunk_codec"
	"example.com/distributed-timeseries-engine/internal/label_index"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/wal"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Repository struct {
	mu                 sync.RWMutex
	metrics            map[string]metric_domain.Metric
	samples            map[string][]metric_domain.Point
	index              *label_index.Index
	wal                *wal.Log
	blocks             *block_store.Store
	maxSamples         int
	writes, duplicates uint64
}

func Open(dir string, maxSamples int) (*Repository, error) {
	w, e := wal.Open(dir + "/wal.log")
	if e != nil {
		return nil, WrapOpenFailure(dir+"/wal.log", e)
	}
	b, e := block_store.Open(dir + "/blocks")
	if e != nil {
		w.Close()
		return nil, WrapOpenFailure(dir+"/blocks", e)
	}
	r := &Repository{metrics: map[string]metric_domain.Metric{}, samples: map[string][]metric_domain.Point{}, index: label_index.New(), wal: w, blocks: b, maxSamples: maxSamples}
	e = w.Replay(func(x wal.Record) error {
		if x.Kind != "sample" {
			return nil
		}
		var s metric_domain.Sample
		if e := jsonUnmarshal(x.Payload, &s); e != nil {
			return e
		}
		r.apply(s)
		return nil
	})
	if e != nil {
		return nil, WrapOpenFailure(dir, e)
	}
	return r, nil
}
func (r *Repository) apply(s metric_domain.Sample) {
	id := s.SeriesID()
	metricKey := s.Tenant + "/" + s.Metric
	if _, ok := r.metrics[metricKey]; !ok {
		r.metrics[metricKey] = metric_domain.Metric{Tenant: s.Tenant, Name: s.Metric, Unit: s.Unit, Labels: s.Labels.Clone(), CreatedAt: time.Now()}
	}
	r.index.Add(id, s.Labels)
	r.samples[id] = append(r.samples[id], metric_domain.Point{Timestamp: s.Timestamp.UnixNano(), Value: s.Value, Quality: s.Quality})
	r.writes++
}
func (r *Repository) Register(m metric_domain.Metric) error {
	if e := m.Validate(); e != nil {
		return e
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[m.Tenant+"/"+m.Name] = m
	r.index.Add(m.SeriesID(), m.Labels)
	return nil
}
func (r *Repository) Ingest(ctx context.Context, in []metric_domain.Sample) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(in) > r.maxSamples {
		return 0, fmt.Errorf("batch exceeds quota")
	}
	accepted := 0
	for i, s := range in {
		if e := s.Validate(); e != nil {
			return accepted, metric_domain.IngestError{Index: i, Err: e}
		}
		metricKey := s.Tenant + "/" + s.Metric
		if _, ok := r.metrics[metricKey]; !ok {
			r.metrics[metricKey] = metric_domain.Metric{Tenant: s.Tenant, Name: s.Metric, Unit: s.Unit, Labels: s.Labels.Clone(), CreatedAt: time.Now()}
		}
		r.index.Add(s.SeriesID(), s.Labels)
		if err := ctx.Err(); err != nil {
			return accepted, err
		}
		id := s.SeriesID()
		dup := false
		for _, p := range r.samples[id] {
			if p.Timestamp == s.Timestamp.UnixNano() {
				dup = true
				break
			}
		}
		if dup {
			r.duplicates++
			continue
		}
		b, _ := jsonMarshal(s)
		if e := r.wal.Append(wal.Record{Kind: "sample", Payload: b}, true); e != nil {
			return accepted, WrapIngestFailure(s.SeriesID(), e)
		}
		r.apply(s)
		accepted++
	}
	return accepted, nil
}
func (r *Repository) Query(ctx context.Context, metric string, labels []metric_domain.LabelMatcher, start, end time.Time) (map[string][]metric_domain.Point, error) {
	if err := ctx.Err(); err != nil {
		return nil, WrapQueryFailure(metric, err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.index.Match(labels)
	if len(ids) == 0 {
		for id, m := range r.metrics {
			if m.Name == metric {
				ids = append(ids, id)
			}
		}
	}
	out := map[string][]metric_domain.Point{}
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return nil, WrapQueryFailure(metric, err)
		}
		if metric != "" && !strings.Contains(id, "/"+metric+"{") {
			continue
		}
		pts := r.samples[id]
		for _, p := range pts {
			if p.Timestamp >= start.UnixNano() && p.Timestamp < end.UnixNano() {
				out[id] = append(out[id], metric_domain.Point{Timestamp: p.Timestamp, Value: p.Value, Quality: p.Quality})
			}
		}
	}
	return out, nil
}

// QuerySeries implements the query executor reader boundary.
func (r *Repository) QuerySeries(ctx context.Context, metric string, labels []metric_domain.LabelMatcher, start, end time.Time) (map[string][]metric_domain.Point, error) {
	return r.Query(ctx, metric, labels, start, end)
}
func (r *Repository) SealAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, pts := range r.samples {
		if len(pts) == 0 {
			continue
		}
		cp := make([]chunk_codec.Point, len(pts))
		for i, p := range pts {
			cp[i] = chunk_codec.Point{Timestamp: p.Timestamp, Value: p.Value, Quality: uint8(p.Quality)}
		}
		if _, e := r.blocks.Seal(id, cp); e != nil {
			return WrapSealFailure(id, e)
		}
		r.samples[id] = nil
	}
	return nil
}
func (r *Repository) Usage() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, v := range r.samples {
		total += len(v)
	}
	return map[string]any{"metrics": len(r.metrics), "active_samples": total, "writes": r.writes, "duplicates": r.duplicates, "cardinality": r.index.Cardinality(), "blocks": r.blocks.Health()["blocks"]}
}
func (r *Repository) Labels() []string                 { return r.index.Names() }
func (r *Repository) LabelValues(name string) []string { return r.index.Values(name) }
func (r *Repository) Close() error                     { return r.wal.Close() }

// Small indirections keep the storage boundary easy to replace with protobuf codecs.
func jsonMarshal(v any) ([]byte, error)   { return encodingJSONMarshal(v) }
func jsonUnmarshal(b []byte, v any) error { return encodingJSONUnmarshal(b, v) }
