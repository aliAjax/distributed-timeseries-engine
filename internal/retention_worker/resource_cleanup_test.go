package retention_worker

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

type trackedCloser struct {
	active *int
}

func (c trackedCloser) Close() error {
	*c.active--
	return nil
}

func TestDteR06BatchZ06(t *testing.T) {
	active := 0
	acquire := func(id string) (io.Closer, error) {
		if active != 0 {
			return nil, fmt.Errorf("resource still active before %s", id)
		}
		active++
		return trackedCloser{active: &active}, nil
	}
	if err := ProcessRetentionBatch([]string{"a", "b", "c"}, acquire, func(string) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if active != 0 {
		t.Fatalf("active=%d", active)
	}
}

func TestDteR06CleanupZ06(t *testing.T) {
	primary := errors.New("delete failed")
	cleanupErr := errors.New("cleanup failed")
	err := RunWithCleanup(func() error { return primary }, func() error { return cleanupErr })
	if !errors.Is(err, primary) || !errors.Is(err, cleanupErr) {
		t.Fatalf("combined error lost identity: %v", err)
	}
}

func TestDteR06RollbackZ06(t *testing.T) {
	primary := errors.New("mutation failed")
	rollbackErr := errors.New("rollback failed")
	rolledBack := false
	err := RunWithRollback(func() error { return primary }, func() error {
		rolledBack = true
		return rollbackErr
	})
	if !rolledBack || !errors.Is(err, primary) || !errors.Is(err, rollbackErr) {
		t.Fatalf("rolledBack=%v err=%v", rolledBack, err)
	}
}

func TestDteR06CloseZ06(t *testing.T) {
	primary := errors.New("retention failed")
	closeErr := errors.New("close failed")
	err := ResolveCloseResult(primary, closeErr)
	if !errors.Is(err, primary) || !errors.Is(err, closeErr) {
		t.Fatalf("close result lost error identity: %v", err)
	}
}
