package toolkitfx_test

import (
	"testing"

	"github.com/ecumenos-social/toolkit/types"
	"github.com/ecumenos-social/toolkitfx"
	"github.com/stretchr/testify/assert"
)

func TestAppConfig_Validate(t *testing.T) {
	tests := map[string]struct {
		cfg    *toolkitfx.AppConfig
		isErr  bool
		errMsg string
	}{
		"should be ok": {
			cfg: &toolkitfx.AppConfig{
				Name:        "test-network-warden",
				Description: "test-network-warden is warden for network",
				RateLimit:   &types.RateLimit{},
			},
		},
		"should returns error for invalid name": {
			cfg: &toolkitfx.AppConfig{
				Name:        "t",
				Description: "test-network-warden is warden for network",
				RateLimit:   &types.RateLimit{},
			},
			isErr:  true,
			errMsg: "invalid name (requirements: 3 < name <= 50, actual: t)",
		},
		"should returns error for invalid description": {
			cfg: &toolkitfx.AppConfig{
				Name:        "test-network-warden",
				Description: "test",
				RateLimit:   &types.RateLimit{},
			},
			isErr:  true,
			errMsg: "invalid description (requirements: 10 < description <= 1024, actual: test)",
		},
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			err := data.cfg.Validate()
			if data.isErr && err == nil {
				t.Fail()
				return
			}
			if !data.isErr && err == nil {
				return
			}
			if !data.isErr && err != nil {
				t.Error(err.Error())
			}
			assert.Equal(t, data.errMsg, err.Error())
		})
	}
}
