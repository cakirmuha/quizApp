package service

import (
	"github.com/labstack/echo"
)

type options struct {
	db bool

	ec *echo.Echo
}

type Option func(*options)

func WithEcho(e *echo.Echo) Option {
	return func(o *options) {
		o.ec = e
	}
}

func WithDB(toggle bool) Option {
	return func(o *options) {
		o.db = toggle
	}
}
