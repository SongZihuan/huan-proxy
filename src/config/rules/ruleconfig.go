package rules

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/api"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/dir"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/file"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/redirect"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/remotetrust"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/respheader"
	"github.com/SongZihuan/huan-proxy/src/config/rules/match"
)

const (
	ProxyTypeFile     = "file"
	ProxyTypeDir      = "dir"
	ProxyTypeAPI      = "api"
	ProxyTypeRedirect = "redirect"
)

type RuleConfig struct {
	Type string `yaml:"type"`

	match.MatchConfig             `yaml:",inline"`
	remotetrust.RemoteTrustConfig `yaml:",inline"`

	File       file.RuleFileConfig            `yaml:"file"`
	Dir        dir.RuleDirConfig              `yaml:"dir"`
	Api        api.RuleAPIConfig              `yaml:"api"`
	Redirect   redirect.RuleRedirectConfig    `yaml:"redirect"`
	RespHeader respheader.SetRespHeaderConfig `yaml:"response-header"`
}

func (p *RuleConfig) SetDefault() {
	p.MatchConfig.SetDefault()
	p.RemoteTrustConfig.SetDefault()

	if p.Type == ProxyTypeFile {
		p.File.SetDefault()
	} else if p.Type == ProxyTypeDir {
		p.Dir.SetDefault()
	} else if p.Type == ProxyTypeAPI {
		p.Api.SetDefault()
	} else if p.Type == ProxyTypeRedirect {
		p.Redirect.SetDefault()
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
		err := p.File.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else if p.Type == ProxyTypeDir {
		err := p.Dir.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else if p.Type == ProxyTypeAPI {
		err := p.Api.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else if p.Type == ProxyTypeRedirect {
		err := p.Redirect.Check()
		if err != nil && err.IsError() {
			return err
		}
	} else {
		return configerr.NewConfigError("proxy type must be file or dir or api")
	}

	return nil
}
