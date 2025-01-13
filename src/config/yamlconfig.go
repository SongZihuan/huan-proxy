package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlConfig struct {
	GlobalConfig         `yaml:",inline"`
	Http                 HttpConfig `yaml:"http"`
	rules.RuleListConfig `yaml:",inline"`
}

func (y *YamlConfig) init() error {
	return nil
}

func (y *YamlConfig) SetDefault() {
	y.GlobalConfig.SetDefault()
	y.Http.SetDefault(&y.GlobalConfig)
	y.RuleListConfig.SetDefault()
}

func (y *YamlConfig) Check() (err configerr.ConfigError) {
	err = y.GlobalConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.Http.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.RuleListConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (y *YamlConfig) parser() configerr.ParserError {
	file, err := os.ReadFile(flagparser.ConfigFile())
	if err != nil {
		return configerr.NewParserError(err, err.Error())
	}

	err = yaml.Unmarshal(file, y)
	if err != nil {
		return configerr.NewParserError(err, err.Error())
	}

	return nil
}
