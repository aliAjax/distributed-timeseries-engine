package config

import (
	"errors"
	"fmt"
)

var ErrConfigParse = errors.New("config parse failed")

type ConfigParseError struct {
	Line int
	Err  error
}

func (e ConfigParseError) Error() string {
	return fmt.Sprintf("parse config line %d: %v", e.Line, e.Err)
}
func (e ConfigParseError) Unwrap() error { return e.Err }

func WrapConfigParse(line int, err error) error {
	if err == nil {
		return nil
	}
	return ConfigParseError{Line: line, Err: err}
}
