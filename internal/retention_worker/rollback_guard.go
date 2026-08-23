package retention_worker

// RunWithRollback runs the operation; on failure it invokes rollback.
// The primary error is never dropped: when rollback also fails, both errors
// are reported so the original failure is never masked.
func RunWithRollback(run func() error, rollback func() error) error {
	primary := run()
	if primary == nil {
		return nil
	}
	rollbackErr := rollback()
	return ResolveCloseResult(primary, rollbackErr)
}
