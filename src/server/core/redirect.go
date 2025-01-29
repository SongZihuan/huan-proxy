package core

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"github.com/SongZihuan/huan-proxy/src/server/context"
	"net/url"
)

func (c *CoreServer) redirectServer(ctx *context.Context) {
	target := c.redirectRewrite(ctx.Rule.Redirect.Address, ctx.Rule.Redirect.Rewrite)

	if _, err := url.Parse(target); err != nil {
		c.abortServerError(ctx)
		return
	}

	c.statusRedirect(ctx, target, ctx.Rule.Redirect.Code)
}

func (c *CoreServer) redirectRewrite(address string, rewrite *rewritecompile.RewriteCompileConfig) string {
	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(address, rewrite.Target)
	}

	return address
}
