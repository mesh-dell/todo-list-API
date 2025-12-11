package errors

import "errors"

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
