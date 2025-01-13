package api

import (
	"fmt"
	resource "github.com/SongZihuan/huan-proxy"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/rewrite"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/url"
)

type RuleAPIConfig struct {
	Address       string                `yaml:"address"`
	AddPrefixPath string                `yaml:"addpath"`
	SubPrefixPath string                `yaml:"subpath"`
	Rewrite       rewrite.RewriteConfig `yaml:"rewrite"`
	HeaderSet     []*HeaderConfig       `yaml:"headerset"`
	HeaderAdd     []*HeaderConfig       `yaml:"headeradd"`
	HeaderDel     []*HeaderDelConfig    `yaml:"headerdel"`
	QuerySet      []*QueryConfig        `yaml:"queryset"`
	QueryAdd      []*QueryConfig        `yaml:"queryadd"`
	QueryDel      []*QueryDelConfig     `yaml:"querydel"`
	Via           string                `yaml:"via"`
}

func (r *RuleAPIConfig) SetDefault() {
	r.AddPrefixPath = utils.ProcessURLPath(r.AddPrefixPath)
	r.SubPrefixPath = utils.ProcessURLPath(r.SubPrefixPath)

	r.Rewrite.SetDefault()

	for _, h := range r.HeaderSet {
		h.SetDefault()
	}

	for _, h := range r.HeaderAdd {
		h.SetDefault()
	}

	for _, h := range r.HeaderDel {
		h.SetDefault()
	}

	for _, q := range r.QuerySet {
		q.SetDefault()
	}

	for _, q := range r.QueryAdd {
		q.SetDefault()
	}

	for _, q := range r.QueryDel {
		q.SetDefault()
	}

	if r.Via == "" {
		r.Via = resource.Via
	}
}

func (r *RuleAPIConfig) Check() configerr.ConfigError {
	targetURL, err := url.Parse(r.Address)
	if err != nil {
		return configerr.NewConfigError(fmt.Sprintf("Failed to parse target URL: %v", err))
	}

	if targetURL.Opaque != "" {
		return configerr.NewConfigError("proxy address should not have Opaque")
	}

	if targetURL.Path == "/" || targetURL.RawPath == "/" {
		targetURL.Path = ""
		targetURL.RawPath = ""
	}

	if targetURL.Path != "" || targetURL.RawPath != "" {
		return configerr.NewConfigError("proxy address should not have path")
	}

	if targetURL.RawQuery != "" {
		return configerr.NewConfigError("proxy address should not have query")
	}

	if targetURL.User != nil {
		return configerr.NewConfigError("proxy address should not have user information")
	}

	if targetURL.Fragment != "" || targetURL.RawFragment != "" {
		return configerr.NewConfigError("proxy address should not have fragment")
	}

	cfgErr := r.Rewrite.Check()
	if cfgErr != nil && cfgErr.IsError() {
		return cfgErr
	}

	return nil
}
