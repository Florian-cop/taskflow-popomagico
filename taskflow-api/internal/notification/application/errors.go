package application

import "errors"

var (
	ErrForbidden      = errors.New("forbidden")
	ErrChannelUnknown = errors.New("channel unknown")
)
