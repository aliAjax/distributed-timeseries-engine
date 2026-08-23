package config

import (
	"errors"
	"fmt"
)

var ErrConfigValidation = errors.New("config validation failed")

type ConfigValidationError struct {
	Field string
	Err   error
}

func (e ConfigValidationError) Error() string {
	return fmt.Sprintf("validate config field %s: %v", e.Field, e.Err)
}
func (e ConfigValidationError) Unwrap() error { return e.Err }

func WrapConfigValidation(field string, err error) error {
	if err == nil {
		return nil
	}
	return ConfigValidationError{Field: field, Err: err}
}
