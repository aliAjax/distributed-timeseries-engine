package recovery_gate

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"example.com/distributed-timeseries-engine/internal/block_store"
	"example.com/distributed-timeseries-engine/internal/chunk_codec"
	"example.com/distributed-timeseries-engine/internal/repository"
	"example.com/distributed-timeseries-engine/internal/wal"
)

func TestDteR14ChunkZ14(t *testing.T) {
	_, err := chunk_codec.Decode([]byte{1, 2, 3})
	var failure chunk_codec.DecodeFailure
	if !errors.As(err, &failure) || !errors.Is(err, chunk_codec.ErrDecode) || !errors.Is(err, chunk_codec.ErrShortChunk) {
		t.Fatalf("decode identity lost: %v", err)
	}
}

func TestDteR14WalZ14(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal.log")
	if err := os.WriteFile(path, make([]byte, 8), 0600); err != nil {
		t.Fatal(err)
	}
	log, err := wal.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	err = log.Replay(func(wal.Record) error { return nil })
	var failure wal.ReplayFailure
	if !errors.As(err, &failure) || !errors.Is(err, wal.ErrReplay) || !errors.Is(err, wal.ErrCorrupt) {
		t.Fatalf("replay identity lost: %v", err)
	}
}

func TestDteR14BlockZ14(t *testing.T) {
	store, err := block_store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Seal("tenant/metric", nil)
	var failure block_store.BlockWriteFailure
	if !errors.As(err, &failure) || !errors.Is(err, block_store.ErrBlockWrite) || !errors.Is(err, block_store.ErrEmptyBlock) {
		t.Fatalf("block identity lost: %v", err)
	}
}

func TestDteR14RepositoryZ14(t *testing.T) {
	dir := t.TempDir()
	var header [8]byte
	binary.BigEndian.PutUint32(header[:4], 64<<20+1)
	if err := os.WriteFile(filepath.Join(dir, "wal.log"), header[:], 0600); err != nil {
		t.Fatal(err)
	}
	_, err := repository.Open(dir, 100)
	var failure repository.RecoveryFailure
	if !errors.As(err, &failure) || !errors.Is(err, repository.ErrRecovery) || !errors.Is(err, wal.ErrReplay) || !errors.Is(err, wal.ErrCorrupt) {
		t.Fatalf("recovery identity lost: %v", err)
	}
}
