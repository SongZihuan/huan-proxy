package redirectcompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/redirect"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"net/url"
)

type RuleRedirectCompileConfig struct {
	Address   string
	TargetURL *url.URL
	Rewrite   *rewritecompile.RewriteCompileConfig
	Code      int
}

func NewRuleAPICompileConfig(r *redirect.RuleRedirectConfig) (*RuleRedirectCompileConfig, error) {
	rewrite, err := rewritecompile.NewRewriteCompileConfig(&r.Rewrite)
	if err != nil {
		return nil, err
	}

	targetURL, err := url.Parse(r.Address)
	if err != nil {
		return nil, err
	}

	return &RuleRedirectCompileConfig{
		Address:   r.Address,
		TargetURL: targetURL,
		Rewrite:   rewrite,
		Code:      r.Code,
	}, nil
}
