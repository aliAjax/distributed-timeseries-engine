package replication

import "context"

func ConsumeReplicaResults(ctx context.Context, results <-chan error) (int, error) {
	completed := 0
	for {
		select {
		case <-ctx.Done():
			return completed, ctx.Err()
		case err := <-results:
			if err != nil {
				return completed, err
			}
			completed++
		}
	}
}
