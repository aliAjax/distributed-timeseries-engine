package retention_worker

import (
	"context"
	"sync"
	"time"
)

type Policy struct {
	Retain      time.Duration
	DeleteDelay time.Duration
}
type Block struct {
	ID          string
	Created     time.Time
	TombstoneAt time.Time
}
type Worker struct {
	mu     sync.Mutex
	blocks map[string]Block
	policy Policy
}

func New(p Policy) *Worker    { return &Worker{blocks: map[string]Block{}, policy: p} }
func (w *Worker) Add(b Block) { w.mu.Lock(); defer w.mu.Unlock(); w.blocks[b.ID] = b }
func (w *Worker) Sweep(now time.Time) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out []string
	for id, b := range w.blocks {
		if b.TombstoneAt.IsZero() && now.Sub(b.Created) > w.policy.Retain {
			b.TombstoneAt = now.Add(w.policy.DeleteDelay)
			w.blocks[id] = b
		}
		if !b.TombstoneAt.IsZero() && !now.Before(b.TombstoneAt) {
			out = append(out, id)
			delete(w.blocks, id)
		}
	}
	return out
}
func (w *Worker) Run(ctx context.Context, interval time.Duration, remove func(string)) {
	if interval <= 0 {
		interval = time.Minute
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			for _, id := range w.Sweep(now) {
				remove(id)
			}
		}
	}
}
