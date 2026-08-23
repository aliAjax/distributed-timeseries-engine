package retention_worker

func RunWithCleanup(run func() error, cleanup func() error) (err error) {
	defer func() {
		err = cleanup()
	}()
	return run()
}
