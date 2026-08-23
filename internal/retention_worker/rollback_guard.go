package retention_worker

import "errors"

func RunWithRollback(run func() error, rollback func() error) error {
	primary := run()
	if primary == nil {
		return nil
	}
	rollbackErr := rollback()
	if rollbackErr != nil {
		return errors.Join(primary, rollbackErr)
	}
	return primary
}
