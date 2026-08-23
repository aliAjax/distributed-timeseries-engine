package retention_worker

// RunWithCleanup runs the operation and always invokes cleanup afterwards.
// The run error is preserved: cleanup cannot mask a run failure as success.
// When both fail, both errors are reported so the original failure is never lost.
func RunWithCleanup(run func() error, cleanup func() error) (err error) {
	primary := run()
	cleanupErr := cleanup()
	return ResolveCloseResult(primary, cleanupErr)
}
