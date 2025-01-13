package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlConfig struct {
	GlobalConfig         `yaml:",inline"`
	Http                 HttpConfig `yaml:"http"`
	rules.RuleListConfig `yaml:",inline"`
}

func (y *YamlConfig) Init() error {
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

func (y *YamlConfig) Parser(filepath string) configerr.ParserError {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return configerr.NewParserError(err, err.Error())
	}

	err = yaml.Unmarshal(file, y)
	if err != nil {
		return configerr.NewParserError(err, err.Error())
	}

	return nil
}
