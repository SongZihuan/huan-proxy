package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlConfig struct {
	GlobalConfig `yaml:",inline"`
	Http         HttpConfig      `yaml:"http"`
	Rules        ProxyRuleConfig `yaml:"rules"`
}

func (y *YamlConfig) init() error {
	return nil
}

func (y *YamlConfig) setDefault() {
	y.GlobalConfig.setDefault()
	y.Http.setDefault(&y.GlobalConfig)
	y.Rules.setDefault()
	fmt.Printf("TAG DDCC [%s]\n", y.Rules.Rules[0].BasePath)
}

func (y *YamlConfig) check(co *CorsOrigin, ps *ProxyServerConfig) (err ConfigError) {
	err = y.GlobalConfig.check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.Http.check(co)
	if err != nil && err.IsError() {
		return err
	}

	err = y.Rules.check(ps)
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (y *YamlConfig) parser() ParserError {
	file, err := os.ReadFile(flagparser.ConfigFile())
	if err != nil {
		return NewParserError(err, err.Error())
	}

	err = yaml.Unmarshal(file, y)
	if err != nil {
		return NewParserError(err, err.Error())
	}

	return nil
}
