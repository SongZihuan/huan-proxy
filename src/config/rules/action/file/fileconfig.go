package file

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/cors"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type RuleFileConfig struct {
	Path string          `yaml:"path"`
	Cors cors.CorsConfig `yaml:"cors"`
}

func (r *RuleFileConfig) SetDefault() {
	r.Cors.SetDefault()
}

func (r *RuleFileConfig) Check() configerr.ConfigError {
	if r.Path == "" {
		return configerr.NewConfigError("file is empty")
	}

	if !utils.IsFile(r.Path) {
		return configerr.NewConfigError("file is not exists")
	}

	err := r.Cors.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}
