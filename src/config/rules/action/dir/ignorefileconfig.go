package dir

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type IgnoreFileConfig struct {
	Regex utils.StringBool `yaml:"regex"`
	File  string           `yaml:"file"`
}

func (i *IgnoreFileConfig) SetDefault() {
	i.Regex.SetDefaultDisable()
}

func (i *IgnoreFileConfig) Check() configerr.ConfigError {
	if i.File == "" {
		return configerr.NewConfigError("file is empty")
	}

	return nil
}
