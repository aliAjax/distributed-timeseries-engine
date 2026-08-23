package retention_worker

import "io"

// ProcessRetentionBatch acquires a handle for each block, removes it, and
// releases the handle before moving to the next block. Releasing eagerly (not
// via defer inside the loop) bounds the number of simultaneously open handles
// to one, so a large retention batch cannot exhaust file descriptors.
func ProcessRetentionBatch(ids []string, acquire func(string) (io.Closer, error), remove func(string) error) error {
	for _, id := range ids {
		resource, err := acquire(id)
		if err != nil {
			return err
		}
		removeErr := remove(id)
		closeErr := resource.Close()
		if err := ResolveCloseResult(removeErr, closeErr); err != nil {
			return err
		}
	}
	return nil
}
