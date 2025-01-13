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
	Address          string                `yaml:"address"`
	ApiAddPrefixPath string                `yaml:"addprefixpath"` // Api前缀避免重名（yaml键忽略）
	ApiSubPrefixPath string                `yaml:"subprefixpath"` // Api前缀避免重名（yaml键忽略）
	ApiRewrite       rewrite.RewriteConfig `yaml:"rewrite"`
	HeaderSet        []HeaderConfig        `yaml:"headerset"`
	HeaderAdd        []HeaderConfig        `yaml:"headeradd"`
	HeaderDel        []HeaderDelConfig     `yaml:"headerdel"`
	QuerySet         []QueryConfig         `yaml:"queryset"`
	QueryAdd         []QueryConfig         `yaml:"queryadd"`
	QueryDel         []QueryDelConfig      `yaml:"querydel"`
	Via              string                `yaml:"via"`
}

func (r *RuleAPIConfig) SetDefault() {
	r.ApiAddPrefixPath = utils.ProcessPath(r.ApiAddPrefixPath)
	r.ApiSubPrefixPath = utils.ProcessPath(r.ApiSubPrefixPath)

	r.ApiRewrite.SetDefault()

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
	_, err := url.Parse(r.Address)
	if err != nil {
		return configerr.NewConfigError(fmt.Sprintf("Failed to parse target URL: %v", err))
	}

	cfgErr := r.ApiRewrite.Check()
	if cfgErr != nil && cfgErr.IsError() {
		return cfgErr
	}

	return nil
}
