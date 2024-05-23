package toolkitfx

import (
	"fmt"

	idgenerator "github.com/ecumenos-social/id-generator"
)

type ServiceName string

type ServiceVersion string

type AppConfig struct {
	Name        string // string (min: 3 chars, max: 50 chars)
	Description string // string (min:10, max: 1024)
	AddrSuffix  string // string (min: 1 char, max: 10 chars)
	NodeID      int64
	RateLimit   float64
	PDNCapacity int64
	NNCapacity  int64
}

func (cfg *AppConfig) Validate() error {
	if len(cfg.Name) < 3 || len(cfg.Name) > 50 {
		return fmt.Errorf("invalid name (requirements: 3 < name <= 50, actual: %s)", cfg.Name)
	}
	if len(cfg.Description) < 10 || len(cfg.Description) > 1024 {
		return fmt.Errorf("invalid description (requirements: 10 < description <= 1024, actual: %s)", cfg.Description)
	}
	if len(cfg.AddrSuffix) < 1 || len(cfg.AddrSuffix) > 10 {
		return fmt.Errorf("invalid address suffix (requirements: 1 < address_suffix <= 10, actual: %s)", cfg.AddrSuffix)
	}
	if maxNodeValue := int64(1<<idgenerator.TopLevelMachineBits - 1); cfg.NodeID < 0 || cfg.NodeID > maxNodeValue {
		return fmt.Errorf("invalid node id (requirements: 0 < node_id <= %d, actual: %d)", maxNodeValue, cfg.NodeID)
	}

	return nil
}
