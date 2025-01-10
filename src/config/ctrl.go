package config

import (
	"github.com/SongZihuan/huan-proxy/src/flagparser"
)

func newConfig() ConfigStruct {
	return ConfigStruct{
		configReady:   false,
		yamlHasParser: false,
	}
}

func InitConfig() ConfigError {
	if !flagparser.IsReady() {
		return NewConfigError("flag not ready")
	}

	config = newConfig()
	err := config.ready()
	if err != nil && err.IsError() {
		return err
	}

	if !config.configReady {
		return NewConfigError("config not ready")
	}

	return nil
}

func IsReady() bool {
	return config.yamlHasParser && config.configReady
}

func Config() *ConfigStruct {
	if !IsReady() {
		panic("config not ready")
	}

	tmp := config
	return &tmp
}

var config ConfigStruct
