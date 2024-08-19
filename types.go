package toolkitfx

import (
	"fmt"

	"github.com/ecumenos-social/toolkit/types"
)

type ServiceName string

type ServiceVersion string

type GenericAppConfig struct {
	ID          int64
	IDGenNode   int64
	Name        string // string (min: 3 chars, max: 50 chars)
	Description string // string (min:10, max: 1024)
	RateLimit   *types.RateLimit
}

func (cfg *GenericAppConfig) Validate() error {
	if len(cfg.Name) < 3 || len(cfg.Name) > 50 {
		return fmt.Errorf("invalid name (requirements: 3 < name <= 50, actual: %s)", cfg.Name)
	}
	if len(cfg.Description) < 10 || len(cfg.Description) > 1024 {
		return fmt.Errorf("invalid description (requirements: 10 < description <= 1024, actual: %s)", cfg.Description)
	}

	return nil
}

type NetworkWardenAppConfig struct {
	AddressSuffix string // string (min:1, max: 30)
}

func (cfg *NetworkWardenAppConfig) Validate() error {
	if len(cfg.AddressSuffix) < 1 || len(cfg.AddressSuffix) > 30 {
		return fmt.Errorf("invalid address suffix (requirements: 1 < address suffix <= 30, actual: %s)", cfg.AddressSuffix)
	}

	return nil
}

type PersonalDataNodeAppConfig struct {
	NetworkWardenID  int64
	Label            string // string (min:1, max: 30)
	AccountsCapacity int64
	CrawlRateLimit   *types.RateLimit
}

func (cfg *PersonalDataNodeAppConfig) Validate() error {
	if len(cfg.Label) < 1 || len(cfg.Label) > 30 {
		return fmt.Errorf("invalid label (requirements: 1 < label <= 30, actual: %s)", cfg.Label)
	}

	return nil
}

type NetworkNodeAppConfig struct {
	NetworkWardenID  int64
	DomainName       string // string (min:1, max: 30)
	AccountsCapacity int64
	CrawlRateLimit   *types.RateLimit
}

func (cfg *NetworkNodeAppConfig) Validate() error {
	if len(cfg.DomainName) < 1 || len(cfg.DomainName) > 30 {
		return fmt.Errorf("invalid domain name (requirements: 1 < domain name <= 30, actual: %s)", cfg.DomainName)
	}

	return nil
}
