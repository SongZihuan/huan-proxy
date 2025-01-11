package config

import (
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
}

func (y *YamlConfig) check(co *CorsOrigin, ps *ProxyServerConfig, ifile *IndexFileCompileList, igfile *IgnoreFileCompileList) (err ConfigError) {
	err = y.GlobalConfig.check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.Http.check(co)
	if err != nil && err.IsError() {
		return err
	}

	err = y.Rules.check(ps, ifile, igfile)
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
