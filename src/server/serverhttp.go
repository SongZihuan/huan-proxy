package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"net/http"
)

func (s *HTTPServer) NormalServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.writeHuanProxyHeader(r)

	func() {
		for _, rule := range s.GetRulesList() {
			if !s.matchURL(rule, r) {
				continue
			}

			if !s.checkProxyTrust(rule, w, r) {
				return
			}

			fmt.Printf("rule.Type: %d\n", rule.Type)

			if rule.Type == rulescompile.ProxyTypeFile {
				s.fileServer(rule, w, r)
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
