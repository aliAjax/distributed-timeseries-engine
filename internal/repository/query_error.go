package repository

import (
	"errors"
	"fmt"
)

var ErrStorageQuery = errors.New("storage query failed")

type QueryFailure struct {
	Metric  string
	Message string
}

func (e QueryFailure) Error() string { return fmt.Sprintf("query %s: %s", e.Metric, e.Message) }

func WrapQueryFailure(metric string, err error) error {
	if err == nil {
		return nil
	}
	return QueryFailure{Metric: metric, Message: err.Error()}
}
