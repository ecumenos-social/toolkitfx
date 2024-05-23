package toolkitfx

import (
	"fmt"
)

type ServiceName string

type ServiceVersion string

type AppConfig struct {
	Name        string // string (min: 3 chars, max: 50 chars)
	Description string // string (min:10, max: 1024)
	RateLimit   float64
}

func (cfg *AppConfig) Validate() error {
	if len(cfg.Name) < 3 || len(cfg.Name) > 50 {
		return fmt.Errorf("invalid name (requirements: 3 < name <= 50, actual: %s)", cfg.Name)
	}
	if len(cfg.Description) < 10 || len(cfg.Description) > 1024 {
		return fmt.Errorf("invalid description (requirements: 10 < description <= 1024, actual: %s)", cfg.Description)
	}

	return nil
}
