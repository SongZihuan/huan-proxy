package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"os"
)

type ConfigStruct struct {
	configReady   bool
	yamlHasParser bool
	sigChan       chan os.Signal

	Yaml  YamlConfig
	Rules *rulescompile.RuleListCompileConfig
}

func (c *ConfigStruct) Parser() configerr.ParserError {
	err := c.Yaml.parser()
	if err != nil {
		return err
	}

	c.yamlHasParser = true
	return nil
}

func (c *ConfigStruct) SetDefault() {
	if !c.yamlHasParser {
		panic("yaml must parser first")
	}

	c.Yaml.SetDefault()
}

func (c *ConfigStruct) Check() (err configerr.ConfigError) {
	err = c.Yaml.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (c *ConfigStruct) CompileRule() configerr.ConfigError {
	res, err := rulescompile.NewRuleListConfig(&c.Yaml.RuleListConfig)
	if err != nil {
		return configerr.NewConfigError("compile rule error: " + err.Error())
	}

	c.Rules = res
	return nil
}

func (c *ConfigStruct) Init() (err configerr.ConfigError) {
	if c.configReady {
		return c.Reload()
	}

	initErr := c.init()
	if initErr != nil {
		return configerr.NewConfigError("init error: " + initErr.Error())
	}

	parserErr := c.Parser()
	if parserErr != nil {
		return configerr.NewConfigError("parser error: " + parserErr.Error())
	} else if !c.yamlHasParser {
		return configerr.NewConfigError("parser error: unknown")
	}

	c.SetDefault()

	err = c.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = c.CompileRule()
	if err != nil && err.IsError() {
		return err
	}

	c.configReady = true
	return nil
}

func (c *ConfigStruct) Reload() (err configerr.ConfigError) {
	if !c.configReady {
		return c.Init()
	}

	bak := *c

	defer func() {
		if err != nil {
			*c = bak
		}
	}()

	reloadErr := c.reload()
	if reloadErr != nil {
		return configerr.NewConfigError("reload error: " + reloadErr.Error())
	}

	parserErr := c.Parser()
	if parserErr != nil {
		return configerr.NewConfigError("reload parser error: " + parserErr.Error())
	} else if !c.yamlHasParser {
		return configerr.NewConfigError("reload parser error: unknown")
	}

	c.SetDefault()

	err = c.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = c.CompileRule()
	if err != nil && err.IsError() {
		return err
	}

	c.configReady = true
	return nil
}

func (c *ConfigStruct) clear() error {
	c.configReady = false
	c.yamlHasParser = false
	// sigChan 不变
	c.Yaml = YamlConfig{}
	return nil
}

func (c *ConfigStruct) init() error {
	c.configReady = false
	c.yamlHasParser = false

	c.sigChan = make(chan os.Signal)
	err := initSignal(c.sigChan)
	if err != nil {
		return err
	}

	err = c.Yaml.init()
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigStruct) reload() error {
	err := c.clear()
	if err != nil {
		return err
	}

	err = c.Yaml.init()
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigStruct) GetSignalChan() chan os.Signal {
	return c.sigChan
}
