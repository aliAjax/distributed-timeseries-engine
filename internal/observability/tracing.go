package observability

import (
	"context"
	"log"
	"sync"
	"time"
)

type TraceID string
type Span struct {
	Name       string
	Trace      TraceID
	Started    time.Time
	Ended      time.Time
	Attributes map[string]string
	Err        error
}
type Tracer struct {
	mu     sync.Mutex
	spans  []Span
	logger *log.Logger
	max    int
}

func NewTracer(logger *log.Logger) *Tracer {
	if logger == nil {
		logger = log.Default()
	}
	return &Tracer{logger: logger, max: 1000, spans: make([]Span, 0)}
}
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	id, ok := ctx.Value(traceKey{}).(TraceID)
	if !ok {
		id = TraceID(time.Now().UTC().Format("20060102T150405.000000000"))
	}
	s := &Span{Name: name, Trace: id, Started: time.Now(), Attributes: map[string]string{}}
	return context.WithValue(ctx, traceKey{}, id), s
}
func (t *Tracer) End(s *Span, err error) {
	if s == nil {
		return
	}
	s.Ended = time.Now()
	s.Err = err
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.spans) >= t.max {
		copy(t.spans, t.spans[1:])
		t.spans = t.spans[:t.max-1]
	}
	t.spans = append(t.spans, snapshotSpan(s))
	t.logger.Printf("trace=%s span=%s duration=%s err=%v", s.Trace, s.Name, s.Ended.Sub(s.Started), err)
}
func (t *Tracer) Recent() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Span, len(t.spans))
	for i := range t.spans {
		out[i] = snapshotSpan(&t.spans[i])
	}
	return out
}

type traceKey struct{}

func WithAttribute(s *Span, key, value string) {
	if s != nil {
		if s.Attributes == nil {
			s.Attributes = map[string]string{}
		}
		s.Attributes[key] = value
	}
}

// snapshotSpan returns a defensive copy of s whose Attributes map is
// independent of s, so later mutations to the live span (or to a returned
// copy) cannot bleed into stored history.
func snapshotSpan(s *Span) Span {
	cp := Span{
		Name:    s.Name,
		Trace:   s.Trace,
		Started: s.Started,
		Ended:   s.Ended,
		Err:     s.Err,
	}
	if s.Attributes != nil {
		attrs := make(map[string]string, len(s.Attributes))
		for k, v := range s.Attributes {
			attrs[k] = v
		}
		cp.Attributes = attrs
	}
	return cp
}

type Timer struct {
	started time.Time
	stop    func(time.Duration)
}

func StartTimer(stop func(time.Duration)) Timer { return Timer{started: time.Now(), stop: stop} }
func (t Timer) Stop() {
	if t.stop != nil {
		t.stop(time.Since(t.started))
	}
}
