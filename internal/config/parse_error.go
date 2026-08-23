package config

import (
	"errors"
	"fmt"
)

var ErrConfigParse = errors.New("config parse failed")

type ConfigParseError struct {
	Line    int
	Message string
}

func (e ConfigParseError) Error() string {
	return fmt.Sprintf("parse config line %d: %s", e.Line, e.Message)
}

func WrapConfigParse(line int, err error) error {
	if err == nil {
		return nil
	}
	return ConfigParseError{Line: line, Message: err.Error()}
}
