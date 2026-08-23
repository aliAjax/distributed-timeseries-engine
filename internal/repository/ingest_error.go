package repository

import (
	"errors"
	"fmt"
)

var ErrStorageIngest = errors.New("storage ingest failed")

type IngestFailure struct {
	Series  string
	Message string
}

func (e IngestFailure) Error() string { return fmt.Sprintf("ingest %s: %s", e.Series, e.Message) }

func WrapIngestFailure(series string, err error) error {
	if err == nil {
		return nil
	}
	return IngestFailure{Series: series, Message: err.Error()}
}
