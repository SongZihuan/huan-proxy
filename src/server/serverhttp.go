package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"net/http"
)

func (s *HuanProxyServer) NormalServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.writeHuanProxyHeader(w, r)

	func() {
		for _, rule := range s.GetRulesList() {
			if !s.matchURL(rule, r) {
				continue
			}

			ctx := NewContext(rule, w, r)

			if !s.checkProxyTrust(ctx) {
				return
			}

			if rule.Type == rulescompile.ProxyTypeFile {
				s.fileServer(ctx)
				return
			} else if rule.Type == rulescompile.ProxyTypeDir {
				s.dirServer(rule, w, r)
				return
			} else if rule.Type == rulescompile.ProxyTypeAPI {
				s.apiServer(rule, w, r)
				return
			} else if rule.Type == rulescompile.ProxyTypeRedirect {
				s.redirectServer(rule, w, r)
				return
			} else {
				s.abortServerError(w)
				return
			}
		}

		s.abortNotFound(w)
	}()
}
