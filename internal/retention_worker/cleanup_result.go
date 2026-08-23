package retention_worker

import "errors"

func RunWithCleanup(run func() error, cleanup func() error) error {
	primary := run()
	cleanupErr := cleanup()
	if primary != nil && cleanupErr != nil {
		return errors.Join(primary, cleanupErr)
	}
	if primary != nil {
		return primary
	}
	return cleanupErr
}
