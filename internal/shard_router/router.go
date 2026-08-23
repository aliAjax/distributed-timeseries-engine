package shard_router

import (
	"hash/fnv"
	"sync"
)

type Assignment struct {
	Shard    int      `json:"shard"`
	Epoch    uint64   `json:"epoch"`
	Replicas []string `json:"replicas"`
}
type Router struct {
	mu          sync.RWMutex
	shards      int
	assignments map[int]Assignment
}

func New(n int) *Router {
	if n < 1 {
		n = 1
	}
	r := &Router{shards: n, assignments: map[int]Assignment{}}
	for i := 0; i < n; i++ {
		r.assignments[i] = Assignment{Shard: i, Epoch: 1, Replicas: []string{"local"}}
	}
	return r
}
func (r *Router) Route(key string) Assignment {
	h := fnv.New32a()
	h.Write([]byte(key))
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.assignments[int(h.Sum32())%r.shards]
}
func (r *Router) List() []Assignment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Assignment, 0, len(r.assignments))
	for _, a := range r.assignments {
		out = append(out, a)
	}
	return out
}
func (r *Router) Assign(a Assignment) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	old := r.assignments[a.Shard]
	if a.Epoch <= old.Epoch {
		return false
	}
	r.assignments[a.Shard] = a
	return true
}
