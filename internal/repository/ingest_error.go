package repository

import (
	"errors"
	"fmt"
)

var ErrStorageIngest = errors.New("storage ingest failed")

type IngestFailure struct {
	Series string
	Err    error
}

func (e IngestFailure) Error() string { return fmt.Sprintf("ingest %s: %v", e.Series, e.Err) }
func (e IngestFailure) Unwrap() error { return e.Err }

func WrapIngestFailure(series string, err error) error {
	if err == nil {
		return nil
	}
	return IngestFailure{Series: series, Err: err}
}
