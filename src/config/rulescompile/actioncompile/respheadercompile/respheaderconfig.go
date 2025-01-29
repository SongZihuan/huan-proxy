package respheadercompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/respheader"
)

type SetRespHeaderCompileConfig struct {
	HeaderSet []*RespHeaderCompileConfig    `yaml:"headerret"`
	HeaderAdd []*RespHeaderCompileConfig    `yaml:"headeradd"`
	HeaderDel []*RespHeaderDelCompileConfig `yaml:"headerdel"`
}

func NewSetRespHeaderCompileConfig(r *respheader.SetRespHeaderConfig) (*SetRespHeaderCompileConfig, error) {
	HeaderSet := make([]*RespHeaderCompileConfig, 0, len(r.HeaderSet))
	for _, v := range r.HeaderSet {
		h, err := NewRespHeaderCompileConfig(v)
		if err != nil {
			return nil, err
		}
		HeaderSet = append(HeaderSet, h)
	}

	HeaderAdd := make([]*RespHeaderCompileConfig, 0, len(r.HeaderAdd))
	for _, v := range r.HeaderAdd {
		h, err := NewRespHeaderCompileConfig(v)
		if err != nil {
			return nil, err
		}
		HeaderAdd = append(HeaderAdd, h)
	}

	HeaderDel := make([]*RespHeaderDelCompileConfig, 0, len(r.HeaderDel))
	for _, v := range r.HeaderDel {
		h, err := NewRespHeaderDelCompileConfig(v)
		if err != nil {
			return nil, err
		}
		HeaderDel = append(HeaderDel, h)
	}

	return &SetRespHeaderCompileConfig{
		HeaderSet: HeaderSet,
		HeaderAdd: HeaderAdd,
		HeaderDel: HeaderDel,
	}, nil
}
