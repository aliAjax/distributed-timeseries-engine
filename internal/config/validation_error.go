package config

import (
	"errors"
	"fmt"
)

var ErrConfigValidation = errors.New("config validation failed")

type ConfigValidationError struct {
	Field   string
	Message string
}

func (e ConfigValidationError) Error() string {
	return fmt.Sprintf("validate config field %s: %s", e.Field, e.Message)
}

func WrapConfigValidation(field string, err error) error {
	if err == nil {
		return nil
	}
	return ConfigValidationError{Field: field, Message: err.Error()}
}
