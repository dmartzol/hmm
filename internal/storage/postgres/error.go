package postgres

import "errors"

var (
	ErrResourceDoesNotExist = errors.New("the resource does not exist")
	ErrExpiredResource      = errors.New("the resource is expired")
)
