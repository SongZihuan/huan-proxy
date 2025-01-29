package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type ReqHeaderCompileConfig struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
}

func NewReqHeaderCompileConfig(h *api.ReqHeaderConfig) (*ReqHeaderCompileConfig, error) {
	return &ReqHeaderCompileConfig{
		Header: h.Header,
		Value:  h.Value,
	}, nil
}
