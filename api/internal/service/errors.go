package service

import "errors"

var (
	ErrDuplicatedUser     = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user haven't registered yet")
	ErrPasswordNotMatched = errors.New("password not matched")
	ErrInternal           = errors.New("something went wrong")
	ErrSubIsEmpty         = errors.New("jwt subject is required")
	ErrInvalidPayload     = errors.New("invalid payload")
)
