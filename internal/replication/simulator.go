package replication

import (
	"context"
	"sync"
	"time"
)

type Replica struct {
	ID      string
	Healthy bool
	Delay   time.Duration
}
type Simulator struct {
	mu       sync.RWMutex
	replicas map[string]Replica
	quorum   int
}

func New(ids []string, quorum int) *Simulator {
	s := &Simulator{replicas: map[string]Replica{}, quorum: quorum}
	for _, id := range ids {
		s.replicas[id] = Replica{ID: id, Healthy: true}
	}
	if quorum < 1 {
		s.quorum = 1
	}
	return s
}
func (s *Simulator) Append(ctx context.Context, fn func(Replica) error) error {
	s.mu.RLock()
	rs := make([]Replica, 0, len(s.replicas))
	for _, r := range s.replicas {
		rs = append(rs, r)
	}
	q := s.quorum
	s.mu.RUnlock()
	ch := make(chan error, len(rs))
	for _, r := range rs {
		go func(x Replica) {
			if x.Delay > 0 {
				select {
				case <-time.After(x.Delay):
				case <-ctx.Done():
					ch <- ctx.Err()
					return
				}
			}
			if !x.Healthy {
				ch <- context.DeadlineExceeded
				return
			}
			ch <- fn(x)
		}(r)
	}
	ok := 0
	for range rs {
		if e := <-ch; e == nil {
			ok++
		}
	}
	if ok < q {
		return context.DeadlineExceeded
	}
	return nil
}
func (s *Simulator) SetHealthy(id string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.replicas[id]
	r.Healthy = ok
	s.replicas[id] = r
}
func (s *Simulator) List() []Replica {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Replica, 0, len(s.replicas))
	for _, r := range s.replicas {
		out = append(out, r)
	}
	return out
}
