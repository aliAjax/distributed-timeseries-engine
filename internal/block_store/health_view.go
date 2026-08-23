package block_store

import "sync"

type HealthView struct {
	mu     sync.RWMutex
	status map[string]int
}

func NewHealthView() *HealthView { return &HealthView{status: make(map[string]int)} }

func (v *HealthView) Set(name string, value int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.status[name] = value
}

func (v *HealthView) Snapshot() map[string]int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make(map[string]int, len(v.status))
	for name, value := range v.status {
		out[name] = value
	}
	return out
}
