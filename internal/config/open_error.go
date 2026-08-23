package config

import (
	"errors"
	"fmt"
)

var ErrConfigOpen = errors.New("config open failed")

type ConfigOpenError struct {
	Path string
	Err  error
}

func (e ConfigOpenError) Error() string { return fmt.Sprintf("open config %s: %v", e.Path, e.Err) }
func (e ConfigOpenError) Unwrap() error { return e.Err }

func WrapConfigOpen(path string, err error) error {
	if err == nil {
		return nil
	}
	return ConfigOpenError{Path: path, Err: err}
}
