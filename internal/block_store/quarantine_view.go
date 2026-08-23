package block_store

import "sync"

type QuarantineView struct {
	mu      sync.RWMutex
	reasons map[string]string
}

func NewQuarantineView() *QuarantineView { return &QuarantineView{reasons: make(map[string]string)} }

func (v *QuarantineView) Put(blockID, reason string) {
	v.reasons[blockID] = reason
}

func (v *QuarantineView) Snapshot() map[string]string {
	return v.reasons
}
