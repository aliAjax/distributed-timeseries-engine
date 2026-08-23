package block_store

import "sync"

type QuarantineView struct {
	mu      sync.RWMutex
	reasons map[string]string
}

func NewQuarantineView() *QuarantineView { return &QuarantineView{reasons: make(map[string]string)} }

func (v *QuarantineView) Put(blockID, reason string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.reasons[blockID] = reason
}

func (v *QuarantineView) Snapshot() map[string]string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make(map[string]string, len(v.reasons))
	for k, val := range v.reasons {
		out[k] = val
	}
	return out
}
