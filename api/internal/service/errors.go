package service

import "errors"

var (
	ErrDuplicatedUser = errors.New("user already exists")
	ErrInternal       = errors.New("something went wrong")
	ErrSubIsEmpty     = errors.New("jwt subject is required")
	ErrInvalidPayload = errors.New("invalid payload")
)
