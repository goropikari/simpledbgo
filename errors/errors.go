package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Err(err error, word string) error {
	return Wrap(err, fmt.Sprintf("failed to %s", word))
}

func Is(err1, err2 error) bool {
	return errors.Is(err1, err2)
}
