package core

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
	"github.com/SongZihuan/huan-proxy/src/server/context"
	"net/http"
)

func (c *CoreServer) cors(corsRule *corscompile.CorsCompileConfig, ctx *context.Context) bool {
	if corsRule.Ignore {
		if ctx.Request.Method() == http.MethodOptions {
			c.abortMethodNotAllowed(ctx)
			return false
		} else {
			return true
		}
	}

	origin := ctx.Request.Header().Get("Origin")
	if origin == "" {
		c.abortForbidden(ctx)
		return false
	}

	if !corsRule.InOriginList(origin) {
		c.abortForbidden(ctx)
		return false
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", corsRule.MaxAgeSec))

	return true
}
