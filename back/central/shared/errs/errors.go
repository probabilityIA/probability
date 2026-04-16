package errs

import (
	"errors"
	"fmt"
)

// ErrWithoutRetry indica que un error no debe reintentarse
var ErrWithoutRetry = errors.New("error without retry")

// Wrap envuelve un error con un mensaje adicional
func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
