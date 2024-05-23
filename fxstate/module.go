package fxstate

import (
	"go.uber.org/fx"
)

var Module = func(initial map[string]interface{}) fx.Option {
	return fx.Options(
		fx.Provide(func() State {
			return New(initial)
		}),
	)
}
