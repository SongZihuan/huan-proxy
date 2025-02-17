package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"os"
)

func InitConfig(configPath string) configerr.ConfigError {
	var err error
	config, err = newConfig(configPath)
	if err != nil {
		return configerr.NewConfigError(err.Error())
	}

	cfgErr := config.Init()
	if cfgErr != nil && cfgErr.IsError() {
		return cfgErr
	}

	if !config.IsReady() {
		return configerr.NewConfigError("config not ready")
	}

	return nil
}

func ReloadConfig() configerr.ConfigError {
	err := config.Reload()
	if err != nil && err.IsError() {
		return err
	}

	if !config.IsReady() {
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

func GetConfigPathFile() string {
	return config.GetConfigPathFile()
}

func GetConfigFileDir() string {
	return config.GetConfigFileDir()
}

func GetConfigFileName() string {
	return config.GetConfigFileName()
}

var config *ConfigStruct
