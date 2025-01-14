package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"os"
)

func InitConfig(configPath string) configerr.ConfigError {
	if !flagparser.IsReady() {
		return configerr.NewConfigError("flag not ready")
	}

	config = newConfig(configPath)
	err := config.Init()
	if err != nil && err.IsError() {
		return err
	}

	if !config.configReady {
		return configerr.NewConfigError("config not ready")
	}

	return nil
}

func IsReady() bool {
	return config.IsReady()
}

func GetConfig() *YamlConfig {
	return config.GetConfig()
}

func GetRules() *rulescompile.RuleListCompileConfig {
	return config.GetRulesList()
}

func GetSignalChan() chan os.Signal {
	return config.GetSignalChan()
}

func NotifyConfigFile() error {
	return config.NotifyConfigFile()
}

func CloseNotifyConfigFile() {
	config.CloseNotifyConfigFile()
}

var config ConfigStruct
