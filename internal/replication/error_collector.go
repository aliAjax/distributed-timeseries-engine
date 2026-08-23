package replication

import "context"

func CollectReplicaError(ctx context.Context, errors chan<- error, err error) bool {
	errors <- err
	return true
}
