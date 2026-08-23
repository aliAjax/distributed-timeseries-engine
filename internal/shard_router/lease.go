package shard_router

import (
	"fmt"
	"sync"
	"time"
)

type Lease struct {
	Shard     int
	Owner     string
	Epoch     uint64
	ExpiresAt time.Time
}
type LeaseTable struct {
	mu     sync.Mutex
	leases map[int]Lease
	ttl    time.Duration
}

func NewLeaseTable(ttl time.Duration) *LeaseTable {
	if ttl <= 0 {
		ttl = time.Minute
	}
	return &LeaseTable{leases: map[int]Lease{}, ttl: ttl}
}
func (l *LeaseTable) Acquire(shard int, owner string, now time.Time) (Lease, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	old, ok := l.leases[shard]
	if ok && now.Before(old.ExpiresAt) && old.Owner != owner {
		return Lease{}, fmt.Errorf("shard lease held by %s", old.Owner)
	}
	x := Lease{Shard: shard, Owner: owner, Epoch: old.Epoch + 1, ExpiresAt: now.Add(l.ttl)}
	l.leases[shard] = x
	return x, nil
}
func (l *LeaseTable) Renew(x Lease, now time.Time) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	old, ok := l.leases[x.Shard]
	if !ok || old.Owner != x.Owner || old.Epoch != x.Epoch {
		return fmt.Errorf("lease fenced")
	}
	if now.After(old.ExpiresAt) {
		return fmt.Errorf("lease expired")
	}
	old.ExpiresAt = now.Add(l.ttl)
	l.leases[x.Shard] = old
	return nil
}
func (l *LeaseTable) List() []Lease {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Lease, 0, len(l.leases))
	for _, x := range l.leases {
		out = append(out, x)
	}
	return out
}
