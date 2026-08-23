package block_store

import "sync"

type CatalogView struct {
	mu      sync.RWMutex
	entries map[string]IndexEntry
}

func NewCatalogView() *CatalogView { return &CatalogView{entries: make(map[string]IndexEntry)} }

func (v *CatalogView) Put(entry IndexEntry) {
	v.entries[entry.BlockID] = entry
}

func (v *CatalogView) Snapshot() map[string]IndexEntry {
	return v.entries
}
