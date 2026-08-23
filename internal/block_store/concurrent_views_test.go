package block_store

import (
	"fmt"
	"sync"
	"testing"
)

type concurrentViewOps struct {
	write func(int)
	read  func() int
}

func exerciseConcurrentView(t *testing.T, ops concurrentViewOps) {
	t.Helper()
	start := make(chan struct{})
	var workers sync.WaitGroup
	workers.Add(2)
	go func() {
		defer workers.Done()
		<-start
		for index := 0; index < 512; index++ {
			ops.write(index)
		}
	}()
	go func() {
		defer workers.Done()
		<-start
		for index := 0; index < 512; index++ {
			_ = ops.read()
		}
	}()
	close(start)
	workers.Wait()
}

func TestDteR03CatalogZ03(t *testing.T) {
	view := NewCatalogView()
	view.Put(IndexEntry{BlockID: "stable", Series: "cpu"})
	exerciseConcurrentView(t, concurrentViewOps{
		write: func(index int) { view.Put(IndexEntry{BlockID: fmt.Sprintf("b-%d", index), Series: "cpu"}) },
		read:  func() int { return len(view.Snapshot()) },
	})
	snapshot := view.Snapshot()
	snapshot["stable"] = IndexEntry{BlockID: "changed"}
	if view.Snapshot()["stable"].BlockID != "stable" {
		t.Fatal("catalog snapshot aliases internal entries")
	}
}

func TestDteR03QuarantineZ03(t *testing.T) {
	view := NewQuarantineView()
	view.Put("stable", "checksum")
	exerciseConcurrentView(t, concurrentViewOps{
		write: func(index int) { view.Put(fmt.Sprintf("b-%d", index), "checksum") },
		read:  func() int { return len(view.Snapshot()) },
	})
	snapshot := view.Snapshot()
	snapshot["stable"] = "changed"
	if view.Snapshot()["stable"] != "checksum" {
		t.Fatal("quarantine snapshot aliases internal reasons")
	}
}

func TestDteR03SeriesZ03(t *testing.T) {
	view := NewSeriesView()
	view.Add("stable", "b-0")
	exerciseConcurrentView(t, concurrentViewOps{
		write: func(index int) { view.Add("cpu", fmt.Sprintf("b-%d", index)) },
		read:  func() int { return len(view.Snapshot()["cpu"]) },
	})
	snapshot := view.Snapshot()
	snapshot["stable"][0] = "changed"
	if view.Snapshot()["stable"][0] != "b-0" {
		t.Fatal("series snapshot aliases internal slices")
	}
}

func TestDteR03HealthZ03(t *testing.T) {
	view := NewHealthView()
	view.Set("stable", 1)
	exerciseConcurrentView(t, concurrentViewOps{
		write: func(index int) { view.Set(fmt.Sprintf("node-%d", index), index) },
		read:  func() int { return len(view.Snapshot()) },
	})
	snapshot := view.Snapshot()
	snapshot["stable"] = 9
	if view.Snapshot()["stable"] != 1 {
		t.Fatal("health snapshot aliases internal status")
	}
}
