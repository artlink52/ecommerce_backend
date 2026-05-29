package errors

import (
	"errors"
)

var (
	ErrProductExists = errors.New("product already exists")
)
