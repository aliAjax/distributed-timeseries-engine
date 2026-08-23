package block_store

import "sync"

type CatalogView struct {
	mu      sync.RWMutex
	entries map[string]IndexEntry
}

func NewCatalogView() *CatalogView { return &CatalogView{entries: make(map[string]IndexEntry)} }

func (v *CatalogView) Put(entry IndexEntry) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.entries[entry.BlockID] = entry
}

func (v *CatalogView) Snapshot() map[string]IndexEntry {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make(map[string]IndexEntry, len(v.entries))
	for k, val := range v.entries {
		out[k] = val
	}
	return out
}
