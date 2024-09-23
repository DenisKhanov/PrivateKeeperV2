package cerrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this login already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
)
