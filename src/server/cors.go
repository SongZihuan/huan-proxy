package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
	"net/http"
)

func (s *HuanProxyServer) cors(corsRule *corscompile.CorsCompileConfig, ctx *Context) bool {
	if corsRule.Ignore {
		if ctx.Request.Method() == http.MethodOptions {
			s.abortMethodNotAllowed(ctx)
			return false
		} else {
			return true
		}
	}

	origin := ctx.Request.Header().Get("Origin")
	if origin == "" {
		s.abortForbidden(ctx)
		return false
	}

	if !corsRule.InOriginList(origin) {
		s.abortForbidden(ctx)
		return false
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", corsRule.MaxAgeSec))

	return true
}
