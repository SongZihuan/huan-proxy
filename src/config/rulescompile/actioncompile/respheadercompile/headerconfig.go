package respheadercompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/respheader"
)

type RespHeaderCompileConfig struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
}

func NewRespHeaderCompileConfig(h *respheader.RespHeaderConfig) (*RespHeaderCompileConfig, error) {
	return &RespHeaderCompileConfig{
		Header: h.Header,
		Value:  h.Value,
	}, nil
}
