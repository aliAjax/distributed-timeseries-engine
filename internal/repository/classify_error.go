package repository

import "errors"

func ClassifyStorageError(err error) string {
	switch {
	case errors.Is(err, ErrStorageOpen):
		return "open"
	case errors.Is(err, ErrStorageIngest):
		return "ingest"
	case errors.Is(err, ErrStorageQuery):
		return "query"
	case errors.Is(err, ErrStorageSeal):
		return "seal"
	default:
		return "unknown"
	}
}
