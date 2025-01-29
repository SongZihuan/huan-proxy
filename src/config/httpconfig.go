package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
)

type HttpConfig struct {
	Address        string `yaml:"address"`
	StopWaitSecond int    `yaml:"stopwaitsecond"`
}

func (h *HttpConfig) SetDefault() {
	if h.Address == "" {
		h.Address = ":4022"
	}

	if h.StopWaitSecond <= 0 {
		h.StopWaitSecond = 10
	}
}

func (h *HttpConfig) Check() configerr.ConfigError {
	return nil
}
