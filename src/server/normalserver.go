package server

import (
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"net/http"
)

func (s *HuanProxyServer) NormalServeHTTP(w http.ResponseWriter, r *http.Request) {
	func() {
	RuleCycle:
		for _, rule := range s.GetRulesList() {
			if !s.matchURL(rule, r) {
				continue RuleCycle
			}

			ctx := NewContext(rule, w, r)

			if !s.checkProxyTrust(ctx) {
				return
			}

			s.writeHuanProxyHeader(ctx)

			if rule.Type == rulescompile.ProxyTypeFile {
				s.fileServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeDir {
				s.dirServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeAPI {
				s.apiServer(ctx)
			} else if rule.Type == rulescompile.ProxyTypeRedirect {
				s.redirectServer(ctx)
			} else {
				s.abortServerError(ctx)
			}

			if config.GetConfig().NotAbort.IsEnable(false) {
				_ = ctx.Reset()
				continue RuleCycle
			}

			ctx.MustWriteToResponse()
			return

		}
		s.defaultResponse(w)
	}()
}
