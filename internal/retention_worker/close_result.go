package retention_worker

import "errors"

// ResolveCloseResult reports the outcome of an operation paired with a close.
// The primary operation error is never dropped: when both fail, both are
// preserved so callers do not record a success when the operation failed.
func ResolveCloseResult(primary, closeErr error) error {
	if primary != nil && closeErr != nil {
		return errors.Join(primary, closeErr)
	}
	if primary != nil {
		return primary
	}
	return closeErr
}
