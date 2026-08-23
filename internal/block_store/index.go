package block_store

import (
	"sort"
	"sync"
	"time"
)

type IndexEntry struct {
	Series     string
	Start, End int64
	BlockID    string
	Points     int
}
type IndexSnapshot struct {
	GeneratedAt time.Time
	Entries     []IndexEntry
}
type BlockIndex struct {
	mu      sync.RWMutex
	entries []IndexEntry
}

func NewIndex() *BlockIndex { return &BlockIndex{entries: make([]IndexEntry, 0)} }
func (i *BlockIndex) Add(e IndexEntry) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.entries = append(i.entries, e)
}
func (i *BlockIndex) Snapshot() IndexSnapshot {
	i.mu.RLock()
	defer i.mu.RUnlock()
	x := append([]IndexEntry(nil), i.entries...)
	return IndexSnapshot{GeneratedAt: time.Now(), Entries: x}
}
func (i *BlockIndex) Find(series string, start, end int64) []IndexEntry {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var out []IndexEntry
	for _, e := range i.entries {
		if e.Series == series && e.End >= start && e.Start < end {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Start < out[b].Start })
	return out
}
func (i *BlockIndex) RemoveBlock(id string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	out := i.entries[:0]
	for _, e := range i.entries {
		if e.BlockID != id {
			out = append(out, e)
		}
	}
	i.entries = out
}
