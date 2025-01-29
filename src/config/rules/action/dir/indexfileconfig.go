package dir

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type IndexFileConfig struct {
	Regex utils.StringBool `yaml:"regex"`
	File  string           `yaml:"file"`
}

func (i *IndexFileConfig) SetDefault() {
	i.Regex.SetDefaultDisable()
}

func (i *IndexFileConfig) Check() configerr.ConfigError {
	if i.File == "" {
		return configerr.NewConfigError("file is empty")
	}

	return nil
}
