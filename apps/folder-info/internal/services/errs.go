package services

import "errors"

var (
	ErrPathIsNotAllowed = errors.New("path-is-not-allowed")
	ErrUnknown          = errors.New("unknown-error")
)
