package repository

import (
	"errors"
	"testing"
)

func TestDteR01OpenZ01(t *testing.T) {
	if !errors.Is(WrapOpenFailure("wal.log", ErrStorageOpen), ErrStorageOpen) {
		t.Fatal("open failure detached from storage sentinel")
	}
}

func TestDteR01IngestZ01(t *testing.T) {
	if !errors.Is(WrapIngestFailure("tenant/cpu", ErrStorageIngest), ErrStorageIngest) {
		t.Fatal("ingest failure detached from storage sentinel")
	}
}

func TestDteR01QueryZ01(t *testing.T) {
	if !errors.Is(WrapQueryFailure("temperature", ErrStorageQuery), ErrStorageQuery) {
		t.Fatal("query failure detached from storage sentinel")
	}
}

func TestDteR01SealZ01(t *testing.T) {
	if !errors.Is(WrapSealFailure("block-7", ErrStorageSeal), ErrStorageSeal) {
		t.Fatal("seal failure detached from storage sentinel")
	}
}
