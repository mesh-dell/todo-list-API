package errors

import "errors"

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrItemNotFound       = errors.New("todo item not found")
	ErrCannotUpdateItem   = errors.New("Forbidden cannot update this todo item")
	ErrCannotDeleteItem   = errors.New("Forbidden cannot delete this todo item")
	ErrCannotGetItem      = errors.New("Forbidden cannot get this todo item")
)
