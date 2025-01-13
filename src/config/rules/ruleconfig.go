package rules

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/api"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/dir"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/file"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/remotetrust"
	"github.com/SongZihuan/huan-proxy/src/config/rules/match"
)

const (
	ProxyTypeFile = "file"
	ProxyTypeDir  = "dir"
	ProxyTypeAPI  = "api"
)

type RuleConfig struct {
	Type string `yaml:"type"`

	match.MatchConfig             `yaml:",inline"`
	remotetrust.RemoteTrustConfig `yaml:",inline"`

	file.RuleFileConfig `yaml:",inline"`
	dir.RuleDirConfig   `yaml:",inline"`
	api.RuleAPIConfig   `yaml:",inline"`
}

func (p *RuleConfig) SetDefault() {
	p.MatchConfig.SetDefault()
	p.RemoteTrustConfig.SetDefault()

	if p.Type == ProxyTypeFile {
		p.RuleFileConfig.SetDefault()
	} else if p.Type == ProxyTypeDir {
		p.RuleDirConfig.SetDefault()
	} else if p.Type == ProxyTypeAPI {
		p.RuleAPIConfig.SetDefault()
	}
}

func (p *RuleConfig) Check() configerr.ConfigError {
	err := p.MatchConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = p.RemoteTrustConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	if p.Type == ProxyTypeFile {
		err := p.RuleFileConfig.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else if p.Type == ProxyTypeDir {
		err := p.RuleDirConfig.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else if p.Type == ProxyTypeAPI {
		err := p.RuleAPIConfig.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else {
		return configerr.NewConfigError("proxy type must be file or dir or api")
	}

	return nil
}
