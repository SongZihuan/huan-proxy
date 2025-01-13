package file

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/cors"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type RuleFileConfig struct {
	File     string          `yaml:"file"`
	FileCors cors.CorsConfig `yaml:"cors"` // File前缀避免重名，（yaml键忽略）
}

func (r *RuleFileConfig) SetDefault() {
	r.FileCors.SetDefault()
}

func (r *RuleFileConfig) Check() configerr.ConfigError {
	if r.File == "" {
		return configerr.NewConfigError("file is empty")
	}

	if utils.IsFile(r.File) {
		return configerr.NewConfigError("file is not exists")
	}

	err := r.FileCors.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}
