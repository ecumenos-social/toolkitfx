package fxidgenerator

import (
	idgenerator "github.com/ecumenos-social/id-generator"
	"go.uber.org/fx"
)

type Config struct {
	TopNodeID int64
	LowNodeID int64
}

var Module = fx.Options(
	fx.Provide(func(config *Config) (idgenerator.Generator, error) {
		return idgenerator.New(&idgenerator.NodeID{
			Top: config.TopNodeID,
			Low: config.LowNodeID,
		})
	}),
)
