package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type HeaderCompileConfig struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
}

func NewHeaderCompileConfig(h *api.HeaderConfig) (*HeaderCompileConfig, error) {
	return &HeaderCompileConfig{
		Header: h.Header,
		Value:  h.Value,
	}, nil
}
