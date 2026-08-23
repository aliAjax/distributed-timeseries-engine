package transport

import (
	"context"
	"errors"
	"testing"
	"time"
)

type requestKey string

func TestDteR09CancelZ09(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	err := BridgeCancellation(ctx, func(context.Context) error {
		called = true
		return nil
	})
	if !errors.Is(err, context.Canceled) || called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}

func TestDteR09DeadlineZ09(t *testing.T) {
	want := time.Now().Add(time.Minute).Round(0)
	ctx, cancel := context.WithDeadline(context.Background(), want)
	defer cancel()
	got, ok := BridgeDeadline(ctx)
	if !ok || !got.Equal(want) {
		t.Fatalf("deadline=%v ok=%v want=%v", got, ok, want)
	}
}

func TestDteR09FreshZ09(t *testing.T) {
	var slot ContextSlot
	first := context.WithValue(context.Background(), requestKey("id"), "first")
	second := context.WithValue(context.Background(), requestKey("id"), "second")
	if got := slot.Load(first).Value(requestKey("id")); got != "first" {
		t.Fatalf("first value=%v", got)
	}
	if got := slot.Load(second).Value(requestKey("id")); got != "second" {
		t.Fatalf("second request inherited value=%v", got)
	}
}

func TestDteR09LoggingZ09(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestKey("tenant"), "plant-a")
	if got := RequestContextForLog(ctx).Value(requestKey("tenant")); got != "plant-a" {
		t.Fatalf("tenant value=%v", got)
	}
}
