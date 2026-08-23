package retention_worker

func ResolveCloseResult(primary, closeErr error) error {
	if primary != nil && closeErr != nil {
		return closeErr
	}
	if primary != nil {
		return primary
	}
	return closeErr
}
