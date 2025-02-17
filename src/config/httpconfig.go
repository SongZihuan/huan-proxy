package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type HttpConfig struct {
	Address        string           `yaml:"address"`
	StopWaitSecond int              `yaml:"stop-wait-second"`
	ProxyProto     utils.StringBool `yaml:"proxy-proto"`
}

func (h *HttpConfig) SetDefault() {
	if h.Address == "" {
		h.Address = ":4022"
	}

	if h.StopWaitSecond <= 0 {
		h.StopWaitSecond = 10
	}

	h.ProxyProto.SetDefaultEnable()
}

func (h *HttpConfig) Check() configerr.ConfigError {
	return nil
}
