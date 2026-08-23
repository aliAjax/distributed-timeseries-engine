package quota

import (
	"fmt"
	"sync"
	"time"
)

type TenantQuota struct {
	MaxSamples int
	MaxSeries  int
	Window     time.Duration
}
type Limiter struct {
	mu    sync.Mutex
	q     map[string]TenantQuota
	used  map[string]int
	reset map[string]time.Time
}

func New() *Limiter {
	return &Limiter{q: map[string]TenantQuota{}, used: map[string]int{}, reset: map[string]time.Time{}}
}
func (l *Limiter) Set(tenant string, q TenantQuota) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.q[tenant] = q
}
func (l *Limiter) Allow(tenant string, n int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	q, ok := l.q[tenant]
	if !ok {
		return nil
	}
	now := time.Now()
	if now.Sub(l.reset[tenant]) >= q.Window {
		l.used[tenant] = 0
		l.reset[tenant] = now
	}
	if q.MaxSamples > 0 && l.used[tenant]+n > q.MaxSamples {
		return fmt.Errorf("tenant %s sample quota exceeded", tenant)
	}
	l.used[tenant] += n
	return nil
}
func (l *Limiter) Usage(tenant string) int { l.mu.Lock(); defer l.mu.Unlock(); return l.used[tenant] }
