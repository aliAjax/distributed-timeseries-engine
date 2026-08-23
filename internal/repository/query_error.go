package repository

import (
	"errors"
	"fmt"
)

var ErrStorageQuery = errors.New("storage query failed")

type QueryFailure struct {
	Metric string
	Err    error
}

func (e QueryFailure) Error() string { return fmt.Sprintf("query %s: %v", e.Metric, e.Err) }
func (e QueryFailure) Unwrap() error { return e.Err }

func WrapQueryFailure(metric string, err error) error {
	if err == nil {
		return nil
	}
	return QueryFailure{Metric: metric, Err: err}
}
