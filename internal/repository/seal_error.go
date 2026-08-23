package repository

import (
	"errors"
	"fmt"
)

var ErrStorageSeal = errors.New("storage seal failed")

type SealFailure struct {
	Block string
	Err   error
}

func (e SealFailure) Error() string { return fmt.Sprintf("seal %s: %v", e.Block, e.Err) }
func (e SealFailure) Unwrap() error { return e.Err }

func WrapSealFailure(block string, err error) error {
	if err == nil {
		return nil
	}
	return SealFailure{Block: block, Err: err}
}
