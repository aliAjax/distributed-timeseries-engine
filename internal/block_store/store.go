package block_store

import (
	"encoding/json"
	"errors"
	"example.com/distributed-timeseries-engine/internal/chunk_codec"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var ErrBlockWrite = errors.New("block write failed")
var ErrEmptyBlock = errors.New("empty block")

type BlockWriteFailure struct {
	Kind error
	Err  error
}

func (e BlockWriteFailure) Error() string { return fmt.Sprintf("block write: %v", e.Err) }

func (e BlockWriteFailure) Unwrap() []error { return []error{e.Kind, e.Err} }

type Block struct {
	ID       string              `json:"id"`
	Series   string              `json:"series"`
	Start    int64               `json:"start"`
	End      int64               `json:"end"`
	Points   []chunk_codec.Point `json:"points"`
	Checksum uint32              `json:"checksum"`
}
type Store struct {
	mu     sync.RWMutex
	dir    string
	blocks map[string]Block
	bad    map[string]string
}

func Open(dir string) (*Store, error) {
	if e := os.MkdirAll(dir, 0750); e != nil {
		return nil, e
	}
	s := &Store{dir: dir, blocks: map[string]Block{}, bad: map[string]string{}}
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	for _, x := range entries {
		if filepath.Ext(x.Name()) != ".json" {
			continue
		}
		b, e := os.ReadFile(filepath.Join(dir, x.Name()))
		if e != nil {
			continue
		}
		var bl Block
		if e = json.Unmarshal(b, &bl); e != nil {
			s.bad[x.Name()] = e.Error()
			continue
		}
		s.blocks[bl.ID] = bl
	}
	return s, nil
}
func (s *Store) Seal(series string, points []chunk_codec.Point) (Block, error) {
	if len(points) == 0 {
		return Block{}, BlockWriteFailure{Kind: ErrBlockWrite, Err: ErrEmptyBlock}
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Timestamp < points[j].Timestamp })
	id := fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(s.blocks))
	bl := Block{ID: id, Series: series, Start: points[0].Timestamp, End: points[len(points)-1].Timestamp, Points: append([]chunk_codec.Point(nil), points...)}
	s.mu.Lock()
	defer s.mu.Unlock()
	b, e := json.MarshalIndent(bl, "", "  ")
	if e != nil {
		return Block{}, e
	}
	if e = os.WriteFile(filepath.Join(s.dir, id+".json"), b, 0600); e != nil {
		return Block{}, e
	}
	s.blocks[id] = bl
	return bl, nil
}
func (s *Store) Query(series string, w metric_domain.TimeWindow) ([]chunk_codec.Point, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []chunk_codec.Point
	for _, b := range s.blocks {
		if b.Series != series {
			continue
		}
		for _, p := range b.Points {
			t := time.Unix(0, p.Timestamp)
			if w.Contains(t) {
				out = append(out, p)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp < out[j].Timestamp })
	return out, nil
}
func (s *Store) List() []Block {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Block, 0, len(s.blocks))
	for _, b := range s.blocks {
		out = append(out, b)
	}
	return out
}
func (s *Store) Quarantine(id, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bad[id] = reason
	delete(s.blocks, id)
	return nil
}
func (s *Store) Health() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]int{"blocks": len(s.blocks), "bad_blocks": len(s.bad)}
}
