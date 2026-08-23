package wal

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestDteR02ReplayZ02(t *testing.T) {
	called := 0
	err := ReplayWithContext(canceledContext(), []Record{{Kind: "sample"}}, func(Record) error {
		called++
		return nil
	})
	if !errors.Is(err, context.Canceled) || called != 0 {
		t.Fatalf("err=%v called=%d", err, called)
	}
}

func TestDteR02AppendZ02(t *testing.T) {
	log, err := Open(filepath.Join(t.TempDir(), "wal.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	err = AppendWithContext(canceledContext(), log, Record{Kind: "sample"}, true)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestDteR02ScanZ02(t *testing.T) {
	visited := 0
	err := ScanWithContext(canceledContext(), []Record{{Kind: "sample"}}, func(int, Record) error {
		visited++
		return nil
	})
	if !errors.Is(err, context.Canceled) || visited != 0 {
		t.Fatalf("err=%v visited=%d", err, visited)
	}
}

func TestDteR02CopyZ02(t *testing.T) {
	var dst bytes.Buffer
	written, err := CopyWithContext(canceledContext(), &dst, bytes.NewBufferString("payload"), make([]byte, 2))
	if !errors.Is(err, context.Canceled) || written != 0 || dst.Len() != 0 {
		t.Fatalf("err=%v written=%d dst=%q", err, written, dst.String())
	}
}
