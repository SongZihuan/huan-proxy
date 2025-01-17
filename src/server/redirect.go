package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"net/http"
	"net/url"
)

func (s *HuanProxyServer) redirectServer(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) {
	target := s.redirectRewrite(rule.Redirect.Address, rule.Redirect.Rewrite)

	if _, err := url.Parse(target); err != nil {
		s.abortServerError(w)
		return
	}

	fmt.Printf("target: %s\n", target)
	s.statusRedirect(w, r, target, rule.Redirect.Code)
}

func (s *HuanProxyServer) redirectRewrite(address string, rewrite *rewritecompile.RewriteCompileConfig) string {
	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(address, rewrite.Target)
	}

	return address
}
