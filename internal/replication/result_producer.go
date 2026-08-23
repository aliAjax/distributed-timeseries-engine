package replication

import (
	"context"
	"sync"
)

func ProduceReplicaResults(ctx context.Context, jobs []func() error) <-chan error {
	results := make(chan error, len(jobs))
	var workers sync.WaitGroup
	workers.Add(len(jobs))
	for _, job := range jobs {
		go func(run func() error) {
			defer workers.Done()
			select {
			case results <- run():
			case <-ctx.Done():
			}
		}(job)
	}
	go func() {
		workers.Wait()
		close(results)
	}()
	return results
}
