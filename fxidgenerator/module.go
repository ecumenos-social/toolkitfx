package fxidgenerator

import (
	"fmt"

	idgenerator "github.com/ecumenos-social/id-generator"
	"go.uber.org/fx"
)

type Config struct {
	TopLevelNode int64
	LowLevelNode int64
}

func (cfg *Config) Validate() error {
	if maxTopLevelNodeValue := int64(1<<idgenerator.TopLevelMachineBits - 1); cfg.TopLevelNode < 0 || cfg.TopLevelNode > maxTopLevelNodeValue {
		return fmt.Errorf("invalid top level node (requirements: 0 < top_level_node <= %d, actual: %d)", maxTopLevelNodeValue, cfg.TopLevelNode)
	}
	if maxLowLevelNodeValue := int64(1<<idgenerator.LowLevelMachineBits - 1); cfg.LowLevelNode < 0 || cfg.LowLevelNode > maxLowLevelNodeValue {
		return fmt.Errorf("invalid top level node (requirements: 0 < top_level_node <= %d, actual: %d)", maxLowLevelNodeValue, cfg.LowLevelNode)
	}
	return nil
}

var Module = fx.Options(
	fx.Provide(func(config *Config) (idgenerator.Generator, error) {
		if err := config.Validate(); err != nil {
			return nil, err
		}

		return idgenerator.New(config.LowLevelNode, config.TopLevelNode)
	}),
)
