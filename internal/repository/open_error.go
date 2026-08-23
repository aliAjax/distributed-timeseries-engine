package repository

import (
	"errors"
	"fmt"
)

var ErrStorageOpen = errors.New("storage open failed")

type OpenFailure struct {
	Path    string
	Message string
}

func (e OpenFailure) Error() string { return fmt.Sprintf("open %s: %s", e.Path, e.Message) }

func WrapOpenFailure(path string, err error) error {
	if err == nil {
		return nil
	}
	return OpenFailure{Path: path, Message: err.Error()}
}
