package replication

import "context"

func ConsumeReplicaResults(ctx context.Context, results <-chan error) (int, error) {
	completed := 0
	for {
		select {
		case <-ctx.Done():
			return completed, ctx.Err()
		case err, ok := <-results:
			if !ok {
				return completed, nil
			}
			if err != nil {
				return completed, err
			}
			completed++
		}
	}
}
