package toolkitfx_test

import (
	"testing"

	"github.com/ecumenos-social/toolkit/types"
	"github.com/ecumenos-social/toolkitfx"
	"github.com/stretchr/testify/assert"
)

func TestNetworkWardenAppConfig_Validate(t *testing.T) {
	tests := map[string]struct {
		cfg    *toolkitfx.NetworkWardenAppConfig
		isErr  bool
		errMsg string
	}{
		"should be ok": {
			cfg: &toolkitfx.NetworkWardenAppConfig{
				Name:          "test-network-warden",
				Description:   "test-network-warden is warden for network",
				AddressSuffix: "nw1",
				RateLimit:     &types.RateLimit{},
			},
		},
		"should returns error for invalid name": {
			cfg: &toolkitfx.NetworkWardenAppConfig{
				Name:          "t",
				Description:   "test-network-warden is warden for network",
				AddressSuffix: "nw1",
				RateLimit:     &types.RateLimit{},
			},
			isErr:  true,
			errMsg: "invalid name (requirements: 3 < name <= 50, actual: t)",
		},
		"should returns error for invalid description": {
			cfg: &toolkitfx.NetworkWardenAppConfig{
				Name:          "test-network-warden",
				Description:   "test",
				AddressSuffix: "nw1",
				RateLimit:     &types.RateLimit{},
			},
			isErr:  true,
			errMsg: "invalid description (requirements: 10 < description <= 1024, actual: test)",
		},
		"should returns error for invalid address suffix": {
			cfg: &toolkitfx.NetworkWardenAppConfig{
				Name:          "test-network-warden",
				Description:   "test-network-warden is warden for network",
				AddressSuffix: "",
				RateLimit:     &types.RateLimit{},
			},
			isErr:  true,
			errMsg: "invalid address suffix (requirements: 1 < address suffix <= 30, actual: )",
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
