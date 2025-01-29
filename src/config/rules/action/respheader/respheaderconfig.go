package respheader

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
)

type SetRespHeaderConfig struct {
	HeaderSet []*RespHeaderConfig    `yaml:"headerret"`
	HeaderAdd []*RespHeaderConfig    `yaml:"headeradd"`
	HeaderDel []*RespHeaderDelConfig `yaml:"headerdel"`
}

func (s *SetRespHeaderConfig) SetDefault() {
	for _, h := range s.HeaderSet {
		h.SetDefault()
	}

	for _, h := range s.HeaderAdd {
		h.SetDefault()
	}

	for _, h := range s.HeaderDel {
		h.SetDefault()
	}
}

func (s *SetRespHeaderConfig) Check() configerr.ConfigError {
	for _, h := range s.HeaderSet {
		err := h.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	for _, h := range s.HeaderAdd {
		err := h.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	for _, h := range s.HeaderDel {
		err := h.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	return nil
}
