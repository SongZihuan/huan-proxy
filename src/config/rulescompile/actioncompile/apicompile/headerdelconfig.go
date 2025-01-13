package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type HeaderDelCompileConfig struct {
	Header string `yaml:"header"`
}

func NewHeaderDelCompileConfig(h *api.HeaderDelConfig) (*HeaderDelCompileConfig, error) {
	return &HeaderDelCompileConfig{
		Header: h.Header,
	}, nil
}
