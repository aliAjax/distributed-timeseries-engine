package repository

import (
	"errors"
	"fmt"
)

var ErrStorageSeal = errors.New("storage seal failed")

type SealFailure struct {
	Block   string
	Message string
}

func (e SealFailure) Error() string { return fmt.Sprintf("seal %s: %s", e.Block, e.Message) }

func WrapSealFailure(block string, err error) error {
	if err == nil {
		return nil
	}
	return SealFailure{Block: block, Message: err.Error()}
}
