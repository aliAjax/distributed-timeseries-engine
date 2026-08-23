package retention_worker

import "errors"

func ResolveCloseResult(primary, closeErr error) error {
	if primary != nil && closeErr != nil {
		return errors.Join(primary, closeErr)
	}
	if primary != nil {
		return primary
	}
	return closeErr
}
