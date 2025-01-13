package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"net/url"
)

type HttpConfig struct {
	Address        string `yaml:"address"`
	StopWaitSecond int    `yaml:"stopwaitsecond"`
}

func (h *HttpConfig) SetDefault(global *GlobalConfig) {
	if h.Address == "" {
		h.Address = "localhost:2689"
	}

	if h.StopWaitSecond <= 0 {
		h.StopWaitSecond = 10
	}
}

func (h *HttpConfig) Check() configerr.ConfigError {
	if _, err := url.Parse(h.Address); err != nil {
		return configerr.NewConfigError(fmt.Sprintf("http address error: %s", err.Error()))
	}
	return nil
}
