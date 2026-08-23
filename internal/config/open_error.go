package config

import (
	"errors"
	"fmt"
)

var ErrConfigOpen = errors.New("config open failed")

type ConfigOpenError struct {
	Path    string
	Message string
}

func (e ConfigOpenError) Error() string { return fmt.Sprintf("open config %s: %s", e.Path, e.Message) }

func WrapConfigOpen(path string, err error) error {
	if err == nil {
		return nil
	}
	return ConfigOpenError{Path: path, Message: err.Error()}
}
