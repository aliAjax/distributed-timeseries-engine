package replication

import "context"

func CollectReplicaError(ctx context.Context, errors chan<- error, err error) bool {
	select {
	case errors <- err:
		return true
	case <-ctx.Done():
		return false
	}
}
