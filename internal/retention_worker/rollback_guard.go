package retention_worker

func RunWithRollback(run func() error, rollback func() error) error {
	primary := run()
	if primary == nil {
		return nil
	}
	rollbackErr := rollback()
	if rollbackErr != nil {
		return rollbackErr
	}
	return primary
}
