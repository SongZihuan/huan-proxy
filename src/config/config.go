package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/fsnotify/fsnotify"
	"os"
	"sync"
)

type ConfigStruct struct {
	ConfigLock sync.Mutex

	configReady   bool
	yamlHasParser bool
	sigchan       chan os.Signal
	configPath    string
	watcher       *fsnotify.Watcher

	Yaml  *YamlConfig
	Rules *rulescompile.RuleListCompileConfig
}

func newConfig(configPath string) ConfigStruct {
	return ConfigStruct{
		// Lock不用初始化
		configReady:   false,
		yamlHasParser: false,
		sigchan:       make(chan os.Signal),
		configPath:    configPath,
		Yaml:          nil,
		Rules:         nil,
	}
}

func (c *ConfigStruct) Init() (err configerr.ConfigError) {
	if c.IsReady() { // 使用IsReady而不是isReady，确保上锁
		return c.Reload()
	}

	initErr := c.init()
	if initErr != nil {
		return configerr.NewConfigError("init error: " + initErr.Error())
	}

	parserErr := c.Parser(c.configPath)
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
	if !c.IsReady() { // 使用IsReady而不是isReady，确保上锁
		return c.Init()
	}

	bak := ConfigStruct{
		configReady:   c.configReady,
		yamlHasParser: c.yamlHasParser,
		sigchan:       c.sigchan,
		configPath:    c.configPath,
		watcher:       c.watcher,
		Yaml:          c.Yaml,
		Rules:         c.Rules,
		// 新建类型
	}

	defer func() {
		if err != nil {
			*c = ConfigStruct{
				configReady:   bak.configReady,
				yamlHasParser: bak.yamlHasParser,
				sigchan:       bak.sigchan,
				configPath:    bak.configPath,
				watcher:       c.watcher,
				Yaml:          bak.Yaml,
				Rules:         bak.Rules,
				// 新建类型 Lock不需要复制
			}
		}
	}()

	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	reloadErr := c.reload()
	if reloadErr != nil {
		return configerr.NewConfigError("reload error: " + reloadErr.Error())
	}

	parserErr := c.Parser(c.configPath)
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
	// sigchan和watcher 不变
	c.Yaml = nil
	c.Rules = nil
	return nil
}

func (c *ConfigStruct) Parser(filepath string) configerr.ParserError {
	err := c.Yaml.Parser(filepath)
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

func (c *ConfigStruct) isReady() bool {
	return c.yamlHasParser && c.configReady
}

func (c *ConfigStruct) init() error {
	c.configReady = false
	c.yamlHasParser = false

	err := initSignal(c.sigchan)
	if err != nil {
		return err
	}

	c.Yaml = new(YamlConfig)
	err = c.Yaml.Init()
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

	c.Yaml = new(YamlConfig)
	err = c.Yaml.Init()
	if err != nil {
		return err
	}

	return nil
}

// export func

func (c *ConfigStruct) IsReady() bool {
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()
	return c.isReady()
}

func (c *ConfigStruct) GetSignalChan() chan os.Signal {
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	return c.sigchan
}

func (c *ConfigStruct) GetConfig() *YamlConfig {
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	if !c.isReady() {
		panic("config is not ready")
	}

	return c.Yaml
}

func (c *ConfigStruct) GetRulesList() *rulescompile.RuleListCompileConfig {
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	if !c.isReady() {
		panic("config is not ready")
	}

	return c.Rules
}
