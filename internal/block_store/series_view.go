package block_store

import "sync"

type SeriesView struct {
	mu     sync.RWMutex
	blocks map[string][]string
}

func NewSeriesView() *SeriesView { return &SeriesView{blocks: make(map[string][]string)} }

func (v *SeriesView) Add(series, blockID string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.blocks[series] = append(v.blocks[series], blockID)
}

func (v *SeriesView) Snapshot() map[string][]string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make(map[string][]string, len(v.blocks))
	for series, blocks := range v.blocks {
		out[series] = append([]string(nil), blocks...)
	}
	return out
}
