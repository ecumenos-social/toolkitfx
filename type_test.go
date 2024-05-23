package toolkitfx_test

import (
	"testing"

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
				AddrSuffix:  "test1",
				NodeID:      7777,
				RateLimit:   0.1,
				PDNCapacity: 10,
				NNCapacity:  10,
			},
		},
		"should returns error for invalid name": {
			cfg: &toolkitfx.AppConfig{
				Name:        "t",
				Description: "test-network-warden is warden for network",
				AddrSuffix:  "test1",
				NodeID:      7777,
				RateLimit:   0.1,
				PDNCapacity: 10,
				NNCapacity:  10,
			},
			isErr:  true,
			errMsg: "invalid name (requirements: 3 < name <= 50, actual: t)",
		},
		"should returns error for invalid description": {
			cfg: &toolkitfx.AppConfig{
				Name:        "test-network-warden",
				Description: "test",
				AddrSuffix:  "test1",
				NodeID:      7777,
				RateLimit:   0.1,
				PDNCapacity: 10,
				NNCapacity:  10,
			},
			isErr:  true,
			errMsg: "invalid description (requirements: 10 < description <= 1024, actual: test)",
		},
		"should returns error for invalid address suffix": {
			cfg: &toolkitfx.AppConfig{
				Name:        "test-network-warden",
				Description: "test-network-warden is warden for network",
				AddrSuffix:  "test110101010101001010101010010101001",
				NodeID:      7777,
				RateLimit:   0.1,
				PDNCapacity: 10,
				NNCapacity:  10,
			},
			isErr:  true,
			errMsg: "invalid address suffix (requirements: 1 < address_suffix <= 10, actual: test110101010101001010101010010101001)",
		},
		"should returns error for invalid node id": {
			cfg: &toolkitfx.AppConfig{
				Name:        "test-network-warden",
				Description: "test-network-warden is warden for network",
				AddrSuffix:  "test1",
				NodeID:      777_777,
				RateLimit:   0.1,
				PDNCapacity: 10,
				NNCapacity:  10,
			},
			isErr:  true,
			errMsg: "invalid node id (requirements: 0 < node_id <= 131071, actual: 777777)",
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
