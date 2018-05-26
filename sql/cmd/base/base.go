package base

import (
	"errors"
)

type DbCmd interface{}

func TypeErrF(msg string) error {
	return errors.New(msg)
}
