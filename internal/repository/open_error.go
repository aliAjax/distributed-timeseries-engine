package repository

import (
	"errors"
	"fmt"
)

var ErrStorageOpen = errors.New("storage open failed")

type OpenFailure struct {
	Path string
	Err  error
}

func (e OpenFailure) Error() string { return fmt.Sprintf("open %s: %v", e.Path, e.Err) }
func (e OpenFailure) Unwrap() error { return e.Err }
func (e OpenFailure) Is(target error) bool { return target == ErrStorageOpen }

func WrapOpenFailure(path string, err error) error {
	if err == nil {
		return nil
	}
	return OpenFailure{Path: path, Err: err}
}
