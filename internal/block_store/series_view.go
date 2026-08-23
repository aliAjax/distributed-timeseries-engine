package block_store

import "sync"

type SeriesView struct {
	mu     sync.RWMutex
	blocks map[string][]string
}

func NewSeriesView() *SeriesView { return &SeriesView{blocks: make(map[string][]string)} }

func (v *SeriesView) Add(series, blockID string) {
	v.blocks[series] = append(v.blocks[series], blockID)
}

func (v *SeriesView) Snapshot() map[string][]string {
	return v.blocks
}
