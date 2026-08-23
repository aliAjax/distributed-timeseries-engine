package block_store

import "sync"

type HealthView struct {
	mu     sync.RWMutex
	status map[string]int
}

func NewHealthView() *HealthView { return &HealthView{status: make(map[string]int)} }

func (v *HealthView) Set(name string, value int) {
	v.status[name] = value
}

func (v *HealthView) Snapshot() map[string]int {
	return v.status
}
