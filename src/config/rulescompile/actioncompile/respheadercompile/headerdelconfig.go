package respheadercompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/respheader"
)

type RespHeaderDelCompileConfig struct {
	Header string `yaml:"header"`
}

func NewRespHeaderDelCompileConfig(h *respheader.RespHeaderDelConfig) (*RespHeaderDelCompileConfig, error) {
	return &RespHeaderDelCompileConfig{
		Header: h.Header,
	}, nil
}
