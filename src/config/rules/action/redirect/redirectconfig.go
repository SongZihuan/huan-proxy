package redirect

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/rewrite"
	"net/http"
	"net/url"
)

type RuleRedirectConfig struct {
	Address string                `yaml:"address"`
	Rewrite rewrite.RewriteConfig `yaml:"rewrite"`
	Code    int                   `yaml:"code"`
}

func (r *RuleRedirectConfig) SetDefault() {
	if r.Code == 0 {
		r.Code = http.StatusMovedPermanently
	}

	r.Rewrite.SetDefault()
}

func (r *RuleRedirectConfig) Check() configerr.ConfigError {
	_, err := url.Parse(r.Address)
	if err != nil {
		return configerr.NewConfigError(fmt.Sprintf("Failed to parse target URL: %s", err.Error()))
	}

	cfgErr := r.Rewrite.Check()
	if cfgErr != nil && cfgErr.IsError() {
		return cfgErr
	}

	if r.Code != http.StatusMovedPermanently && r.Code != http.StatusMultipleChoices {
		return configerr.NewConfigError(fmt.Sprintf("Redirect code must be %d %d: you use %d", http.StatusMovedPermanently, http.StatusMultipleChoices, r.Code))
	}

	return nil
}
