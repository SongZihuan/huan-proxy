package config

import "os"

type ConfigStruct struct {
	configReady   bool
	yamlHasParser bool
	sigChan       chan os.Signal

	Yaml        YamlConfig
	CoreOrigin  CorsOrigin
	ProxyServer ProxyServerConfig
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

	err = c.CoreOrigin.init()
	if err != nil {
		return err
	}

	err = c.ProxyServer.init()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigStruct) parser() ParserError {
	err := c.Yaml.parser()
	if err != nil {
		return err
	}

	c.yamlHasParser = true
	return nil
}

func (c *ConfigStruct) setDefault() {
	if !c.yamlHasParser {
		panic("yaml must parser first")
	}

	c.Yaml.setDefault()
}

func (c *ConfigStruct) check() (err ConfigError) {
	err = c.Yaml.check(&c.CoreOrigin, &c.ProxyServer)
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (c *ConfigStruct) ready() (err ConfigError) {
	if c.configReady {
		return nil
	}

	initErr := c.init()
	if initErr != nil {
		return NewConfigError("init error: " + initErr.Error())
	}

	parserErr := c.parser()
	if parserErr != nil {
		return NewConfigError("parser error: " + parserErr.Error())
	} else if !c.yamlHasParser {
		return NewConfigError("parser error: unknown")
	}

	c.setDefault()
	err = c.check()
	if err != nil && err.IsError() {
		return err
	}

	c.configReady = true
	return nil
}

func (c *ConfigStruct) GetSignalChan() chan os.Signal {
	return c.sigChan
}
