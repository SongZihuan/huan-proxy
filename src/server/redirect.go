package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"net/url"
)

func (s *HuanProxyServer) redirectServer(ctx *Context) {
	target := s.redirectRewrite(ctx.Rule.Redirect.Address, ctx.Rule.Redirect.Rewrite)

	if _, err := url.Parse(target); err != nil {
		s.abortServerError(ctx)
		return
	}

	fmt.Printf("target: %s\n", target)
	s.statusRedirect(ctx, target, ctx.Rule.Redirect.Code)
}

func (s *HuanProxyServer) redirectRewrite(address string, rewrite *rewritecompile.RewriteCompileConfig) string {
	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(address, rewrite.Target)
	}

	return address
}
