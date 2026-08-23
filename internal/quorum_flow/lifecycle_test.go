package quorum_flow_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	replication "example.com/distributed-timeseries-engine/internal/replication"
)

func TestDteR05ProducerZ05(t *testing.T) {
	start := make(chan struct{})
	jobs := []func() error{
		func() error { <-start; return nil },
		func() error { <-start; return errors.New("replica unavailable") },
	}
	results := replication.ProduceReplicaResults(context.Background(), jobs)
	close(start)
	received := 0
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-results:
			if !ok {
				if received != 2 {
					t.Fatalf("received=%d", received)
				}
				return
			}
			received++
		case <-timer.C:
			t.Fatal("replica result channel was not closed")
		}
	}
}

func TestDteR05WaitZ05(t *testing.T) {
	start := make(chan struct{})
	done := make(chan struct{})
	var completed atomic.Int32
	replication.WaitForReplicas(start, []func(){func() { completed.Add(1) }, func() { completed.Add(1) }}, done)
	select {
	case <-done:
		t.Fatal("coordinator completed before workers started")
	case <-time.After(30 * time.Millisecond):
	}
	close(start)
	select {
	case <-done:
		if completed.Load() != 2 {
			t.Fatalf("completed=%d", completed.Load())
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("coordinator did not complete")
	}
}

func TestDteR05ErrorsZ05(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := make(chan struct{})
	done := make(chan bool, 1)
	go func() {
		<-start
		done <- replication.CollectReplicaError(ctx, make(chan error), errors.New("failed"))
	}()
	close(start)
	select {
	case sent := <-done:
		if sent {
			t.Fatal("canceled error send reported success")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("canceled error sender leaked")
	}
}

func TestDteR05ConsumerZ05(t *testing.T) {
	start := make(chan struct{})
	results := make(chan error)
	go func() {
		<-start
		results <- nil
		close(results)
	}()
	type outcome struct {
		count int
		err   error
	}
	done := make(chan outcome, 1)
	go func() {
		<-start
		count, err := replication.ConsumeReplicaResults(context.Background(), results)
		done <- outcome{count: count, err: err}
	}()
	close(start)
	select {
	case result := <-done:
		if result.err != nil || result.count != 1 {
			t.Fatalf("count=%d err=%v", result.count, result.err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("consumer did not exit after channel close")
	}
}
