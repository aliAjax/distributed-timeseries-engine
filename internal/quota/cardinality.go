package quota

import (
	"fmt"
	"sync"
	"time"
)

type CardinalityGuard struct {
	mu   sync.Mutex
	max  int
	seen map[string]time.Time
}

func NewCardinalityGuard(max int) *CardinalityGuard {
	return &CardinalityGuard{max: max, seen: map[string]time.Time{}}
}
func (g *CardinalityGuard) Observe(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.seen[id]; ok {
		return nil
	}
	if g.max > 0 && len(g.seen) >= g.max {
		return fmt.Errorf("cardinality limit reached")
	}
	g.seen[id] = time.Now()
	return nil
}
func (g *CardinalityGuard) Count() int { g.mu.Lock(); defer g.mu.Unlock(); return len(g.seen) }
func (g *CardinalityGuard) Reset()     { g.mu.Lock(); defer g.mu.Unlock(); g.seen = map[string]time.Time{} }
