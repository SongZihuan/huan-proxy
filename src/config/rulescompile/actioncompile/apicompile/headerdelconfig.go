package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type ReqHeaderDelCompileConfig struct {
	Header string `yaml:"header"`
}

func NewReqHeaderDelCompileConfig(h *api.ReqHeaderDelConfig) (*ReqHeaderDelCompileConfig, error) {
	return &ReqHeaderDelCompileConfig{
		Header: h.Header,
	}, nil
}
