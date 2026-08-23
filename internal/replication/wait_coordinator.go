package replication

import "sync"

func WaitForReplicas(start <-chan struct{}, jobs []func(), done chan<- struct{}) {
	var workers sync.WaitGroup
	for _, job := range jobs {
		go func(run func()) {
			defer workers.Done()
			<-start
			workers.Add(1)
			run()
		}(job)
	}
	go func() {
		workers.Wait()
		close(done)
	}()
}
